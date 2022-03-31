################################################################################

# rpmbuilder:relative-pack true

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Simple utility for controlling application working time
Name:            deadline
Version:         1.5.5
Release:         0%{?dist}
Group:           Applications/System
License:         Apache 2.0
URL:             https://kaos.sh/deadline

Source0:         https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.17

Provides:        %{name} = %{version}-%{release}

################################################################################

%description
Simple utility for controlling application working time.

################################################################################

%prep
%setup -q

%build
export GOPATH=$(pwd)
pushd src/github.com/essentialkaos/%{name}
  go build -mod vendor -o $GOPATH/%{name} %{name}.go
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -pm 755 %{name} %{buildroot}%{_bindir}/

%clean
rm -rf %{buildroot}

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_bindir}/%{name}

################################################################################

%changelog
* Wed Mar 30 2022 Anton Novojilov <andy@essentialkaos.com> - 1.5.5-0
- Removed pkg.re usage
- Added module info
- Added Dependabot configuration

* Sun Apr 04 2021 Anton Novojilov <andy@essentialkaos.com> - 1.5.4-0
- Updated compatibility with the latest version of ek

* Fri Dec 04 2020 Anton Novojilov <andy@essentialkaos.com> - 1.5.3-0
- ek package updated to v12

* Wed Dec 04 2019 Anton Novojilov <andy@essentialkaos.com> - 1.5.2-0
- ek package updated to v11

* Thu Dec 13 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.1-0
- Code refactoring

* Fri Nov 02 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.0-0
- Fixed bug with showing version info
- Code refactoring

* Sun May 21 2017 Anton Novojilov <andy@essentialkaos.com> - 1.4.0-0
- ek package updated to v9

* Fri Apr 21 2017 Anton Novojilov <andy@essentialkaos.com> - 1.3.1-0
- Added build tag

* Sun Apr 16 2017 Anton Novojilov <andy@essentialkaos.com> - 1.3.0-0
- ek package updated to v8

* Tue Mar 07 2017 Anton Novojilov <andy@essentialkaos.com> - 1.2.0-0
- ek package updated to latest version

* Fri Jan 13 2017 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Initial build
