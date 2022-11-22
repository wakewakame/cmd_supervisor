FROM debian:11.5

# ユーザーの作成
ARG USERNAME=user
ARG GROUPNAME=user
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID $GROUPNAME && \
    useradd -m -s /bin/bash -u $UID -g $GID $USERNAME

# 必要なパッケージのインストール
RUN apt update -y && apt upgrade -y
RUN apt install -y \
    wget \
    postgresql-client \
    redis-tools

# Golangのインストール
# https://go.dev/doc/install
RUN wget https://go.dev/dl/go1.19.3.linux-amd64.tar.gz -O - | \
    tar -C /usr/local -xzf -

# rootですべき作業が済んだらユーザを切り替える
USER $USERNAME
WORKDIR /home/$USERNAME/workspace

# Golangのパスを通る
RUN echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.bashrc && \
    echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
ENV GO=/usr/local/go/bin/go

# ホットリロードを行うためにairをインストールする
RUN $GO install github.com/cosmtrek/air@v1.40.4

CMD bash
