FROM marpteam/marp-cli:latest

USER root

# 1. Install Python
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    libreoffice \
    && rm -rf /var/lib/apt/lists/*

# 2. Install Marp CLI Globally (THIS FIXES THE PATH ISSUE)
RUN npm install -g @marp-team/marp-cli

# 3. Setup Python Environment
ENV VIRTUAL_ENV=/opt/venv
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
ENV PYTHONUNBUFFERED=1

COPY requirement.txt /tmp/requirement.txt
RUN pip install --no-cache-dir -r /tmp/requirement.txt && rm /tmp/requirement.txt

COPY . /home/marp/app/
RUN chown -R marp:marp /home/marp/app

USER marp
WORKDIR /home/marp/app

# 4. Critical: Reset entrypoint
ENTRYPOINT []

CMD ["python3", "server.py"]