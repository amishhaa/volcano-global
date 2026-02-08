# 1. Base Image - Go 1.24
FROM --platform=linux/arm64 golang:1.24-bullseye

# 2. Install dependencies + Official Docker CLI
RUN apt-get update && apt-get install -y \
    ca-certificates curl gnupg lsb-release && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
    echo "deb [arch=arm64 signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian bullseye stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && apt-get install -y docker-ce-cli && \
    rm -rf /var/lib/apt/lists/*

# 3. Install Kind (ARM64), Kubectl, Helm
RUN curl -Lo /usr/local/bin/kind https://kind.sigs.k8s.io/dl/v0.22.0/kind-linux-arm64 && chmod +x /usr/local/bin/kind
RUN curl -LO "https://dl.k8s.io/release/v1.29.0/bin/linux/arm64/kubectl" && chmod +x kubectl && mv kubectl /usr/local/bin/
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

WORKDIR /workspace
COPY . .

COPY hack/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
