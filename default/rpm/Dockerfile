FROM centos:6

# Go & basic build tools
RUN yum update -y && \
    yum groupinstall -y 'Development Tools' && \
    yum install -y cyrus-sasl-devel openssl-devel libffi readline-devel && \
    curl -LO https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go1.8.3.linux-amd64.tar.gz && \
    rm go1.8.3.linux-amd64.tar.gz

# Install Ruby (for fpm)
RUN git clone https://github.com/rbenv/ruby-build.git && \
    cd ruby-build && \
    ./install.sh && \
    cd .. && \
    rm -rf ruby-build && \
    ruby-build 2.3.1 /usr/local

# Install fpm
RUN gem install fpm
