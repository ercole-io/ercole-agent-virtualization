Name:           ercole-agent-virtualization
Version:        ERCOLE_VERSION
Release:        1%{?dist}
Summary:        Agent Virtualization for ercole

License:        Proprietary
URL:            https://github.com/ercole-io/%{name}
Source0:        https://github.com/ercole-io/%{name}/archive/%{name}-%{version}.tar.gz
Requires: systemd
BuildRequires: systemd

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

install -d $RPM_BUILD_ROOT/opt/ercole-agent-virtualization/run
install -d %{buildroot}%{_unitdir} 
install -d %{buildroot}%{_presetdir}
install -m 0644 package/rhel7/ercole-agent-virtualization.service %{buildroot}%{_unitdir}/%{name}.service
install -m 0644 package/rhel7/60-ercole-agent-virtualization.preset %{buildroot}%{_presetdir}/60-%{name}.preset

%post
/usr/bin/systemctl preset %{name}.service >/dev/null 2>&1 ||:

%preun
/usr/bin/systemctl --no-reload disable %{name}.service >/dev/null 2>&1 || :
/usr/bin/systemctl stop %{name}.service >/dev/null 2>&1 ||:

%postun
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:

%files
%attr(-,ercole,-) /opt/ercole-agent-virtualization/run
%dir /opt/ercole-agent-virtualization
%dir /opt/ercole-agent-virtualization/fetch
%config(noreplace) /opt/ercole-agent-virtualization/config.json
/opt/ercole-agent-virtualization/fetch/filesystem
/opt/ercole-agent-virtualization/fetch/host
/opt/ercole-agent-virtualization/fetch/vmware.ps1
/opt/ercole-agent-virtualization/fetch/ovm
/opt/ercole-agent-virtualization/ercole-agent-virtualization
%{_unitdir}/ercole-agent-virtualization.service
%{_presetdir}/60-ercole-agent-virtualization.preset

%changelog
* Mon May  7 2018 Simone Rota <srota2@sorint.it>
- 
