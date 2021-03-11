.\ximfect.exe about-tool

$amt = 15
$total = 0
$min = 0
$max = 0

for($i = 0; $i -lt $amt; $i++)
{
	Write-Output "--------------------"
	Write-Output ("Benchmark {0}" -f ($i + 1))
	$time = Measure-Command {.\ximfect.exe apply-effect discord_red.png blur discord_red_blur.png}
	if($time.TotalSeconds -gt $max)
	{
		$max = $time.TotalSeconds
	}
	if($min -eq 0)
	{
		$min = $max
	}
	if($time.TotalSeconds -lt $min)
	{
		$min = $time.TotalSeconds
	}
	$total += $time.TotalSeconds
	Write-Output ("Time: {0}" -f $time.TotalSeconds)
}

$avg = $total / $amt
Write-Output "===================="
Write-Output ("Average: {0}" -f $avg) ("Minimum: {0}" -f $min) ("Maximum: {0}" -f $max)
Write-Output "===================="