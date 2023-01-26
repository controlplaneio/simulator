"""
borrowed from jython
https://bitbucket.org/jython/jython/raw/28a66ba038620292520470a0bb4dc9bb8ac2e403/Lib/select.py
"""

import errno
import os

class error(Exception):
    pass

ALL = None

_exception_map = {}

def _map_exception(exc, circumstance=ALL):
    try:
        mapped_exception = _exception_map[(exc.__class__, circumstance)]
        mapped_exception.java_exception = exc
        return mapped_exception
    except KeyError:
        return error(-1, 'Unmapped java exception: <%s:%s>' % (exc.toString(), circumstance))

POLLIN   = 1
POLLOUT  = 2

# The following event types are completely ignored on jython
# Java does not support them, AFAICT
# They are declared only to support code compatibility with cpython

POLLPRI  = 4
POLLERR  = 8
POLLHUP  = 16
POLLNVAL = 32

def _getselectable(selectable_object):
    try:
        channel = selectable_object.getchannel()
    except:
        try:
            channel = selectable_object.fileno().getChannel()
        except:
            raise TypeError("Object '%s' is not watchable" % selectable_object,
                            errno.ENOTSOCK)

    return channel

# Fake Selector
class Selector:

    def close(self):
        pass

    def keys(self):
        return []

    def select(self, timeout=None):
        return []

    def selectedKeys(self):
        class SelectedKeys:
            def iterator(self):
                return []
        return SelectedKeys()

    def selectNow(self, timeout=None):
        return []

class poll:

    def __init__(self):
        self.selector = Selector()
        self.chanmap = {}
        self.unconnected_sockets = []

    def _register_channel(self, socket_object, channel, mask):
        jmask = 0
        if mask & POLLIN:
            # Note that OP_READ is NOT a valid event on server socket channels.
            if channel.validOps() & OP_ACCEPT:
                jmask = OP_ACCEPT
            else:
                jmask = OP_READ
        if mask & POLLOUT:
            if channel.validOps() & OP_WRITE:
                jmask |= OP_WRITE
            if channel.validOps() & OP_CONNECT:
                jmask |= OP_CONNECT
        selectionkey = channel.register(self.selector, jmask)
        self.chanmap[channel] = (socket_object, selectionkey)

    def _check_unconnected_sockets(self):
        temp_list = []
        for socket_object, mask in self.unconnected_sockets:
            channel = _getselectable(socket_object)
            if channel is not None:
                self._register_channel(socket_object, channel, mask)
            else:
                temp_list.append( (socket_object, mask) )
        self.unconnected_sockets = temp_list

    def register(self, socket_object, mask = POLLIN|POLLOUT|POLLPRI):
        try:
            channel = _getselectable(socket_object)
            if channel is None:
                # The socket is not yet connected, and thus has no channel
                # Add it to a pending list, and return
                self.unconnected_sockets.append( (socket_object, mask) )
                return
            self._register_channel(socket_object, channel, mask)
        except BaseException as exc:
            raise _map_exception(exc)

    def unregister(self, socket_object):
        try:
            channel = _getselectable(socket_object)
            self.chanmap[channel][1].cancel()
            del self.chanmap[channel]
        except BaseException as exc:
            raise _map_exception(exc)

    def _dopoll(self, timeout):
        if timeout is None or timeout < 0:
            self.selector.select()
        else:
            try:
                timeout = int(timeout)
                if not timeout:
                    self.selector.selectNow()
                else:
                    # No multiplication required: both cpython and java use millisecond timeouts
                    self.selector.select(timeout)
            except ValueError as vx:
                raise error("poll timeout must be a number of milliseconds or None", errno.EINVAL)
        # The returned selectedKeys cannot be used from multiple threads!
        return self.selector.selectedKeys()

    def poll(self, timeout=None):
        return []

    def _deregister_all(self):
        try:
            for k in self.selector.keys():
                k.cancel()
            # Keys are not actually removed from the selector until the next select operation.
            self.selector.selectNow()
        except BaseException as exc:
            raise _map_exception(exc)

    def close(self):
        try:
            self._deregister_all()
            self.selector.close()
        except BaseException as exc:
            raise _map_exception(exc)

def _calcselecttimeoutvalue(value):
    if value is None:
        return None
    try:
        floatvalue = float(value)
    except Exception as x:
        raise TypeError("Select timeout value must be a number or None")
    if value < 0:
        raise error("Select timeout value cannot be negative", errno.EINVAL)
    if floatvalue < 0.000001:
        return 0
    return int(floatvalue * 1000) # Convert to milliseconds

# This cache for poll objects is required because of a bug in java on MS Windows
# http://bugs.jython.org/issue1291

class poll_object_cache:

    def __init__(self):
        self.is_windows = os.name == 'nt'
        if self.is_windows:
            self.poll_object_queue = Queue.Queue()
        import atexit
        atexit.register(self.finalize)

    def get_poll_object(self):
        if not self.is_windows:
            return poll()
        try:
            return self.poll_object_queue.get(False)
        except Queue.Empty:
            return poll()

    def release_poll_object(self, pobj):
        if self.is_windows:
            pobj._deregister_all()
            self.poll_object_queue.put(pobj)
        else:
            pobj.close()

    def finalize(self):
        if self.is_windows:
            while True:
                try:
                    p = self.poll_object_queue.get(False)
                    p.close()
                except Queue.Empty:
                    return

_poll_object_cache = poll_object_cache()

def native_select(read_fd_list, write_fd_list, outofband_fd_list, timeout=None):
    timeout = _calcselecttimeoutvalue(timeout)
    # First create a poll object to do the actual watching.
    pobj = _poll_object_cache.get_poll_object()
    try:
        registered_for_read = {}
        # Check the read list
        for fd in read_fd_list:
            pobj.register(fd, POLLIN)
            registered_for_read[fd] = 1
        # And now the write list
        for fd in write_fd_list:
            if fd in registered_for_read:
                # registering a second time overwrites the first
                pobj.register(fd, POLLIN|POLLOUT)
            else:
                pobj.register(fd, POLLOUT)
        results = pobj.poll(timeout)
        # Now start preparing the results
        read_ready_list, write_ready_list, oob_ready_list = [], [], []
        for fd, mask in results:
            if mask & POLLIN:
                read_ready_list.append(fd)
            if mask & POLLOUT:
                write_ready_list.append(fd)
        return read_ready_list, write_ready_list, oob_ready_list
    finally:
        _poll_object_cache.release_poll_object(pobj)

select = native_select

def cpython_compatible_select(read_fd_list, write_fd_list, outofband_fd_list, timeout=None):
    # First turn all sockets to non-blocking
    # keeping track of which ones have changed
    modified_channels = []
    try:
        for socket_list in [read_fd_list, write_fd_list, outofband_fd_list]:
            for s in socket_list:
                channel = _getselectable(s)
                if channel.isBlocking():
                    modified_channels.append(channel)
                    channel.configureBlocking(0)
        return native_select(read_fd_list, write_fd_list, outofband_fd_list, timeout)
    finally:
        for channel in modified_channels:
            channel.configureBlocking(1)
