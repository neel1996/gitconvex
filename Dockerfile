FROM golang:1.16.0

WORKDIR /go/src/github.com/neel1996/gitconvex-server

COPY . .

# Install required packages from apt-get
RUN apt-get update && \
    apt-get install apt-transport-https ca-certificates gnupg software-properties-common wget sudo -y && \
    apt-add-repository 'deb https://apt.kitware.com/ubuntu/ focal main' && \
    wget -O - https://apt.kitware.com/keys/kitware-archive-latest.asc 2>/dev/null | gpg --dearmor - | sudo tee /etc/apt/trusted.gpg.d/kitware.gpg >/dev/null && \
    apt-get install cmake -y

# Install node js
RUN curl -fsSL https://deb.nodesource.com/setup_15.x | sudo -E bash - && \
    apt-get install nodejs -y

# Building React UI bundle
RUN cd ui/ && \
    npm install && \
    export NODE_ENV=production && \
    npm install tailwindcss postcss autoprefixer && \
    npx tailwindcss build -o src/index.css -c src/tailwind.config.js && \
    npm run build && \
    mv build/ gitconvex-ui/ && \
    mv gitconvex-ui/ ../ && \
    cd .. && \
    rm -rf ui/

# Download and build OpenSSL
RUN cd ~ && git clone https://github.com/openssl/openssl.git && \
    cd openssl && ./Configure && \
    make && make install && \
    cp -rp *.so* /usr/lib/

# Download libssh2
RUN cd ~ && wget https://github.com/libssh2/libssh2/releases/download/libssh2-1.9.0/libssh2-1.9.0.tar.gz && \
    tar -xzf libssh2-1.9.0.tar.gz && \
    cd libssh2-1.9.0/ && \
    mkdir build && cd build && \
    cmake -DCMAKE_C_COMPILER=gcc -DCMAKE_INSTALL_PREFIX=../install .. && cmake --build . --target install

# Download and Libgit2 setup
RUN cd ~ && wget https://github.com/libgit2/libgit2/releases/download/v1.1.0/libgit2-1.1.0.tar.gz && \
    tar -xzf libgit2-1.1.0.tar.gz && \
    cd libgit2-1.1.0/ && \
    mkdir build && cd build && \
    cmake -DCMAKE_PREFIX_PATH=../../libssh2-1.9.0/install/ -DCMAKE_INSTALL_PREFIX=../install -DBUILD_CLAR=OFF .. && \
    cmake --build . --target install && \
    make install && \
    cp -rp ~/libgit2-1.1.0/install/include/* /usr/include/ && \
    cp -rp ~/libgit2-1.1.0/install/lib/pkgconfig/* /usr/lib/pkgconfig && \
    cp -rp ~/libgit2-1.1.0/install/lib/lib* /usr/lib/

# Post Cleanup stage
RUN apt-get remove cmake nodejs apt-transport-https ca-certificates gnupg software-properties-common wget -y

EXPOSE 9001

CMD export PKG_CONFIG_PATH=/usr/local/lib && \
    go run /go/src/github.com/neel1996/gitconvex-server/server.go