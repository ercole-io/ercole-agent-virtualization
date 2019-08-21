Name:           ercole-agent-virtualization
Version:        ERCOLE_VERSION
Release:        1%{?dist}
Summary:        Agent Virtualization for ercole

License:        Proprietary
URL:            https://github.com/amreo/%{name}
Source0:        https://github.com/amreo/%{name}/archive/%{name}-%{version}.tar.gz

Group:          Tools

Buildroot: /tmp/rpm-ercole-agent-virtualization

%global debug_package %{nil}

%description
Ercole Virtualization Agent collects information about vms and clusters
running on the local machine and send information to a central server

%pre
getent passwd ercole >/dev/null || \
    useradd -r -g oinstall -G oinstall,dba -d /home/ercole-agent-virtualization -m -s /bin/bash \
    -c "Ercole agent user" ercole
getent passwd ercole >/dev/null || \
    useradd -r -g dba -d /home/ercole-agent-virtualization -m -s /bin/bash \
    -c "Ercole agent user" ercole
getent passwd ercole >/dev/null || \
    useradd -r -g oinstall -d /home/ercole-agent-virtualization -m -s /bin/bash \
    -c "Ercole agent user" ercole
exit 0

%prep
%setup -q -n %{name}-%{version}

rm -rf $RPM_BUILD_ROOT
make DESTDIR=$RPM_BUILD_ROOT/opt/ercole-agent-virtualization install
install -d $RPM_BUILD_ROOT/etc/systemd/system
install -d $RPM_BUILD_ROOT/opt/ercole-agent-virtualization/run
install -m 644 package/rhel7/ercole-agent-virtualization.service $RPM_BUILD_ROOT/etc/systemd/system/ercole-agent-virtualization.service

%post

%files
%attr(-,ercole,-) /opt/ercole-agent-virtualization/run
%dir /opt/ercole-agent-virtualization
%dir /opt/ercole-agent-virtualization/fetch
%config(noreplace) /opt/ercole-agent-virtualization/config.json
%config(noreplace) /opt/ercole-agent-virtualization/creds.csv
/opt/ercole-agent-virtualization/fetch/filesystem
/opt/ercole-agent-virtualization/fetch/host
/opt/ercole-agent-virtualization/fetch/vmware.ps1
/opt/ercole-agent-virtualization/ercole-agent-virtualization
/etc/systemd/system/ercole-agent-virtualization.service
%changelog
* Mon May  7 2018 Simone Rota <srota2@sorint.it>
- 
