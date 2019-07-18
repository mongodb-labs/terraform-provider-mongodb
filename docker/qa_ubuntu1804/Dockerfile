FROM ubuntu:18.04
ARG ssh_pub_key
ENV DEBIAN_FRONTEND noninteractive

# Download dependencies
RUN apt-get update && apt-get install -y \
    curl \
    iproute2 \
    libcurl4 \
    net-tools \
    openssh-server \
    sudo \
    ;

# Configure SSHD
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
RUN mkdir -p /var/run/sshd /root/.ssh
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config

# Allow authentication with the specified SSH public key
RUN echo "${ssh_pub_key}" > /root/.ssh/authorized_keys

# Create a data directory for Ops Manager release automation
RUN mkdir -p /data/automation/mongodb-releases

# Expose SSHD and HTTP
EXPOSE 22
EXPOSE 8080

# Start the SSH daemon
CMD ["/usr/sbin/sshd", "-D"]
