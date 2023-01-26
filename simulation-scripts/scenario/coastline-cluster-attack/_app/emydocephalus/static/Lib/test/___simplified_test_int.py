import sys

L = [
        ('0', 0),
        ('1', 1),
        ('9', 9),
        ('10', 10),
        ('99', 99),
        ('100', 100),
        ('314', 314),
        (' 314', 314),
        ('314 ', 314),
        ('  \t\t  314  \t\t  ', 314),
        (repr(sys.maxsize), sys.maxsize),
        ('  1x', ValueError),
        ('  1  ', 1),
        ('  1\02  ', ValueError),
        ('', ValueError),
        (' ', ValueError),
        ('  \t\t  ', ValueError),
        ("\u0200", ValueError)
]

class CtxManager:
    
    def __init__(self, exception, **kw):
        self.exception = exception
        self.kw = kw
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, *args):
        if exc_type is None or not issubclass(exc_type, self.exception):
            raise AssertionError('should have raised %s but raised %s'
                %self.exception, exc_type)
        return True

class TestCase:
    
    def assertEqual(self, result, expected):
        if result != expected:
            raise AssertionError('expected %s, got %s' %(expected, result))

    def assertIsInstance(self, obj, klass):
        if not isinstance(obj, klass):
            print(obj,'is not instance of', klass)

    def assertRaises(self, exception, *args, **kw):
        if args:
            callable, *args = args
            try:
                callable(*args, **kw)
            except Exception as exc:
                if not isinstance(exc, exception):
                    print(callable, args, kw, 'does not raise', exception)
        else:
            return CtxManager(exception, **kw)

    def assertTrue(self, expr):
        if not expr:
            raise AssertionError('expr not True : %s' %expr)

import pickle

class ListTest(TestCase):
    type2test = list

    def test_basic(self):
        self.assertEqual(list([]), [])
        l0_3 = [0, 1, 2, 3]
        l0_3_bis = list(l0_3)
        self.assertEqual(l0_3, l0_3_bis)
        self.assertTrue(l0_3 is not l0_3_bis)
        self.assertEqual(list(()), [])
        self.assertEqual(list((0, 1, 2, 3)), [0, 1, 2, 3])
        self.assertEqual(list(''), [])
        self.assertEqual(list('spam'), ['s', 'p', 'a', 'm'])

        if sys.maxsize == 0x7fffffff:
            # This test can currently only work on 32-bit machines.
            # XXX If/when PySequence_Length() returns a ssize_t, it should be
            # XXX re-enabled.
            # Verify clearing of bug #556025.
            # This assumes that the max data size (sys.maxint) == max
            # address size this also assumes that the address size is at
            # least 4 bytes with 8 byte addresses, the bug is not well
            # tested
            #
            # Note: This test is expected to SEGV under Cygwin 1.3.12 or
            # earlier due to a newlib bug.  See the following mailing list
            # thread for the details:

            #     http://sources.redhat.com/ml/newlib/2002/msg00369.html
            self.assertRaises(MemoryError, list, range(sys.maxsize // 2))

        # This code used to segfault in Py2.4a3
        x = []
        x.extend(-y for y in x)
        self.assertEqual(x, [])

    def test_truth(self):
        super().test_truth()
        self.assertTrue(not [])
        self.assertTrue([42])

    def test_identity(self):
        self.assertTrue([] is not [])

    def test_len(self):
        super().test_len()
        self.assertEqual(len([]), 0)
        self.assertEqual(len([0]), 1)
        self.assertEqual(len([0, 1, 2]), 3)

    def test_overflow(self):
        lst = [4, 5, 6, 7]
        n = int((sys.maxsize*2+2) // len(lst))
        def mul(a, b): return a * b
        def imul(a, b): a *= b
        self.assertRaises((MemoryError, OverflowError), mul, lst, n)
        self.assertRaises((MemoryError, OverflowError), imul, lst, n)

    def test_repr_large(self):
        # Check the repr of large list objects
        def check(n):
            l = [0] * n
            s = repr(l)
            self.assertEqual(s,
                '[' + ', '.join(['0'] * n) + ']')
        check(10)       # check our checking code
        check(1000000)

    def test_iterator_pickle(self):
        # Userlist iterators don't support pickling yet since
        # they are based on generators.
        data = self.type2test([4, 5, 6, 7])
        it = itorg = iter(data)
        print(it)
        d = pickle.dumps(it)
        it = pickle.loads(d)
        self.assertEqual(type(itorg), type(it))
        self.assertEqual(self.type2test(it), self.type2test(data))

        it = pickle.loads(d)
        next(it)
        d = pickle.dumps(it)
        self.assertEqual(self.type2test(it), self.type2test(data)[1:])

    def test_reversed_pickle(self):
        data = self.type2test([4, 5, 6, 7])
        it = itorg = reversed(data)
        d = pickle.dumps(it)
        it = pickle.loads(d)
        self.assertEqual(type(itorg), type(it))
        self.assertEqual(self.type2test(it), self.type2test(reversed(data)))

        it = pickle.loads(d)
        next(it)
        d = pickle.dumps(it)
        self.assertEqual(self.type2test(it), self.type2test(reversed(data))[1:])

    def test_no_comdat_folding(self):
        # Issue 8847: In the PGO build, the MSVC linker's COMDAT folding
        # optimization causes failures in code that relies on distinct
        # function addresses.
        class L(list): pass
        with self.assertRaises(TypeError):
            (3,) + L([1,2])

def run(test_class):
    tester = test_class()
    for attr in dir(tester):
        if attr.startswith('test_'):
            print(attr)
            getattr(tester, attr)()

run(ListTest)


