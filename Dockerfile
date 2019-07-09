FROM centos/tools

ADD . /usr/src/kube-app-netutil

WORKDIR /usr/src/kube-app-netutil

ENV INSTALL_PKGS "golang"
RUN rpm --import https://mirror.go-repo.io/centos/RPM-GPG-KEY-GO-REPO && \
    curl -s https://mirror.go-repo.io/centos/go-repo.repo | tee /etc/yum.repos.d/go-repo.repo && \
    yum install -y $INSTALL_PKGS && \
    rpm -V $INSTALL_PKGS && \
    yum clean all && \
    make clean && \
    make

RUN cp /usr/src/kube-app-netutil/bin/server /usr/bin
RUN cp /usr/src/kube-app-netutil/bin/client /usr/bin

WORKDIR /

LABEL io.k8s.display-name="Kube application netutil"

CMD ["/usr/bin/server"]
