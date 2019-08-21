<#

	.SYNOPSIS
	Automated VMware cluster and VMs available in a list of provided vCenter servers.

	.DESCRIPTION
	Automated VMware cluster and VMs available in a list of provided vCenter servers.

	.EXAMPLE
	./vmware.ps1 -s vms
	
	.EXAMPLE
	./vmware.ps1 -s cluster

	.NOTES
	File Name  : cluster.ps1  
	Author     : Riccardo Suardi - rsuardi@sorint.it 
	Requires   : PowerShell Core, VMware PowerCLI

	.LINK
	https://sorint.it

	.Parameter s
	string, switch variable
	values accepted: cluster or vms
	##

#>

param (
	[Parameter(Mandatory=$true)][string]$s
)
#Set-PowerCLIConfiguration -InvalidCertificateAction:Ignore
if(Test-Path .\creds.csv) {
    $vcs = Import-CSV .\creds.csv
    foreach ($vc in $vcs) {
        Connect-VIServer $vc.server -User $vc.user -Password $vc.pass | Out-Null
        New-VIProperty -Name NumCPU -ObjectType Cluster -Value {
                    $TotalPCPU = 0
                    $Args[0] | Get-VMHost | Foreach {
                        $TotalPCPU += $_.NumCPU
                    }
                    $TotalPCPU
            } `
            -Force -WarningAction:SilentlyContinue | Out-Null
            
        New-VIProperty -Name NumSockets -ObjectType Cluster -Value {
                    $TotalPSOCKS = 0
                    $Args[0] | Get-VMHost | Foreach {
                        $TotalPSOCKS += $_.ExtensionData.Hardware.CpuInfo.NumCpuPackages
                    }
                    $TotalPSOCKS
            } `
            -Force -WarningAction:SilentlyContinue | Out-Null
		switch ($s.ToUpper()) {
			"VMS" {
				# OUTPUT FORMAT: cluster name, vm name, guest os hostname
				Get-VM | Select @{N="Cluster";E={Get-Cluster -VM $_}}, Name, @{N="guestHostname";E={$_.ExtensionData.Guest.HostName}} | ConvertTo-CSV | % { $_ -replace '"', ""}
			}
			"CLUSTER" {
				# OUTPUT FORMAT: cluster name, core sum, socket sum
				Get-Cluster | Select Name, NumCPU, NumSockets | ConvertTo-CSV | % { $_ -replace '"', ""}
			}
			Default	{ Write-Host "wrong switch selection" }
		}
        Disconnect-VIServer $vc.server -Confirm:$false | Out-Null
    }
}
else { Write-Warning "Credentials file not found, please check working dir."}
