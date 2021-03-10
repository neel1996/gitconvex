### Instructions to build libgit2

- [Windows](#for-windows)
- [Linux](#for-linux-tested-on-ubuntu-2004)
- [MacOS](#for-macos)

#### For windows

Make sure you have the following applications installed on your system

- [VS 2019](https://visualstudio.microsoft.com/vs/features/cplusplus/) for the C/C++ compiler
- [cmake](https://cmake.org/download/)
- [cygwin](https://www.cygwin.com/) with the following packages
    - [make](https://cygwin.com/packages/summary/make.html)
    - [pkg-config](https://cygwin.com/packages/summary/pkg-config.html)
    - [libssh2-devel](https://cygwin.com/packages/summary/libssh2-devel.html)
    - [zlib](https://cygwin.com/packages/summary/zlib.html) (optional)

_After setting up cygwin, add cygwin /bin folder to the **path** environment variable_

The libraries required for setting up `libgit2` with `libssh2` and `openssl` are available as a zip file
in [lib/win](lib/win). Extract the content of the zip file into the same folder (lib/win)

**Steps**

- Download the `.zip` file from the libgit2 [releases](https://github.com/libgit2/libgit2/releases/tag/v1.1.0) or clone
  the [repo](https://github.com/libgit2/libgit2) from github

> Note : gitconvex was tested against libgit v1.1.0. So we recommend you to download the zip file from the releases section for stability

Extract / clone the libgit2 source to an easily accessible folder and copy the [lib](lib) folder to the root of libgit2
directory

- Clone [libssh2](https://github.com/libssh2/libssh2) into the root of the libgit2 folder

```shell
git clone https://github.com/libssh2/libssh2 libssh2
```

- Search for `Cygwin64 Terminal` and open it. Run the following commands to generate the libgit2 DLL

```shell
cd <FOLDER WHERE LIBGIT2 IS AVAILABLE>

mkdir build && cd build

cmake -DCMAKE_INSTALL_PREFIX=../install \
      -DCMAKE_FIND_ROOT_PATH=../lib \
      -DOPENSSL_SSL_LIBRARY=../lib/win/lib \
      -DOPENSSL_CRYPTO_LIBRARY=../lib/win/lib \
      -DEMBED_SSH_PATH=../libssh2 ..

# If the above command completes without any error, execute the following command      
cmake --build . --target install
```

- After the command completes execution
    - copy the `git2.dll` from `../bin/` to `<cygwin>/bin` folder
    - copy `git2.lib` from `../lib` to `<cygwin>/lib` folder
    - copy `libgit2.pc` from `../lib/pkgconfig` to `<cygwin>/lib/pkgconfig`

Following the above steps will make `libgit2` available as a shared library for `git2go`.

#### For Linux (tested on Ubuntu 20.04)

Make sure you have the following packages installed on your system

- [cmake](https://cmake.org/download/)
- [curl](https://curl.se/)
- [wget](https://www.gnu.org/software/wget/)
- A suitable C/C++ compiler

Make sure you have required rights for running the following commands. If any of these fail due to access errors, then
try with `sudo`

```shell
# Download and setup openssl
cd ~ && git clone https://github.com/openssl/openssl.git openssl
cd openssl && ./Configure 
make && make install 
cp -rp *.so* /usr/lib/

# Download and setup libssh2
cd ~ && wget https://github.com/libssh2/libssh2/releases/download/libssh2-1.9.0/libssh2-1.9.0.tar.gz
tar -xzf libssh2-1.9.0.tar.gz 
cd libssh2-1.9.0/ 
mkdir build && cd build
cmake -DCMAKE_INSTALL_PREFIX=../install .. && cmake --build . --target install

#Download and setup libgit2
cd ~ && wget https://github.com/libgit2/libgit2/releases/download/v1.1.0/libgit2-1.1.0.tar.gz
tar -xzf libgit2-1.1.0.tar.gz 
cd libgit2-1.1.0/
mkdir build && cd build 
cmake -DCMAKE_PREFIX_PATH=../../libssh2-1.9.0/install/ -DCMAKE_INSTALL_PREFIX=../install -DBUILD_CLAR=OFF .. 
cmake --build . --target install 

# Copy the libgit2 shared object and pkconfig files to the /usr/lib path
cp -rp ~/libgit2-1.1.0/install/include/* /usr/include/ && \
cp -rp ~/libgit2-1.1.0/install/lib/pkgconfig/* /usr/lib/ && \
cp -rp ~/libgit2-1.1.0/install/lib/lib* /usr/lib/
```

#### For MacOS

The steps are similar to the [Linux](#for-linux-tested-on-ubuntu-2004) guidelines mentioned above.

Download and setup [brew](https://brew.sh/) to install all the required packages

- [cmake](https://formulae.brew.sh/formula/cmake)
- [curl](https://formulae.brew.sh/formula/curl)
- [wget](https://formulae.brew.sh/formula/wget)

Once the packages are setup, run the same commands mentioned above for Linux to setup libgit2 