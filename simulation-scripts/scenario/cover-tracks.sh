# clean up fs
find  /etc/kubernetes -exec touch {} +
rm -f ~/.bash_history

rm -f \
  /opt/challenge.txt \
  /root/flag.txt \
  /root/flag-found.txt

# clean up pods
# rm all pods

# clean up deploys etc?

# clean events?
# kubectl delete --raw /api/v1/namespaces/default/ee454724358c72 --help ??

# clean syslogs
find /var/log/journal /run/log/ -name "*.journal" -delete


# clean apt cache
