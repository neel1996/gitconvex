### Instructions to build libgit2

- [Windows](#for-windows)
- [Linux](#for-linux)

#### For windows

Make sure you have the following applications installed on your system

- [cmake](https://cmake.org/download/)
- [cygwin](https://www.cygwin.com/) with the following packages
    - [make](https://cygwin.com/packages/summary/make.html)
    - [pkg-config](https://cygwin.com/packages/summary/pkg-config.html)
    - [zlib](https://cygwin.com/packages/summary/zlib.html) (optional)

_After setting up cygwin, add cygwin /bin folder to the **path** environment variable_

The libraries required for setting up `libgit2` with `libssh2` and `openssl` are available in [lib/win](lib/win)

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

cmake -DCMAKE_INSTALL_PREFIX=../install && \
      -DCMAKE_FIND_ROOT_PATH=../lib && \
      -DOPENSSL_SSL_LIBRARY=../lib/win/lib && \
      -DOPENSSL_CRYPTO_LIBRARY=../lib/win/lib && \
      -DEMBED_SSH_PATH=../libssh2 ..

# If the above command completes without any error, execute the following command      
cmake --build . --target install
```

- After the command completes execution
    - copy the `git2.dll` from `../bin/` to `<cygwin>/bin` folder
    - copy `git2.lib` from `../lib` to `<cygwin>/lib` folder
    - copy `libgit2.pc` from `../lib/pkgconfig` to `<cygwin>/lib/pkgconfig`

Following the above steps will make `libgit2` available as a shared library for `git2go`.