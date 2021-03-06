FROM golang:1.15.5-alpine

WORKDIR /go/src/github.com/neel1996/gitconvex-server

COPY . .

# Alpine package download stage
RUN apk update && apk upgrade
RUN apk add --update nodejs nodejs-npm pkgconfig cmake gcc libc-dev perl git linux-headers make

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
    make && make install

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
    mv ~/libgit2-1.1.0/install/include/* /usr/local/include/ && \
    mv ~/libgit2-1.1.0/install/lib64/pkgconfig/* /usr/local/lib/

RUN apk del nodejs nodejs-npm perl git

EXPOSE 9001

CMD export PKG_CONFIG_PATH=~/libgit2-1.1.0/install/lib64/pkgconfig && \
    go run /go/src/github.com/neel1996/gitconvex-server/server.go
