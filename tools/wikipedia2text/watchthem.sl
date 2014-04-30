#
# watches for php processes that are running too long
#
debug(7);

sub proc
{
   local('$temp');
   $temp = split('\s+', $1);
   if ($temp[0] eq "") { shift($temp); }
   return @($temp[0], $temp[3]);
}

$opid = @();

while (1)
{
   $output = filter({ return iff($1[1] eq "php", $1); }, map(&proc, `ps -a`));
   $pids   = map({ return $1[0]; }, $output);

   $check  = retainAll($opid, $pids);
   map({ warn("Killing $+ $1 $+ '"); `/bin/kill -9 $1`; }, $check);

   $opid   = $pids;
   sleep (120 * 1000);
}
