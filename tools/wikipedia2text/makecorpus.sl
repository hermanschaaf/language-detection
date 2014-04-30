#
# split a file into several corpora files.
#
debug(7 | 34);

# open file handles.

$out = @ARGV[1];

for ($x = 0; $x < 768; $x++)
{
   push(@handles, openf("> $+ $out $+ /split $+ $x $+ .txt"));
}

# ...

$handle = openf(@ARGV[0]);
@collect = @();

while $text (readln($handle))
{
   if ($text ne "")
   {
      push(@collect, $text);
   }
   else
   {
      printAll(rand(@handles), @collect);
      @collect = @();
   }
}

closef($handle);

# close file handles

foreach $handle (@handles)
{
   closef($handle);
}
