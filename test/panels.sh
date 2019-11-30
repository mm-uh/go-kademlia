tmux new-session \; \
  send-keys './node 127.0.0.1 8081 127.0.0.1 8080' C-m \; \
  split-window -v \; \
  send-keys './node 127.0.0.1 8082 127.0.0.1 8080' C-m \; \
  split-window -h \; \
  send-keys './node 127.0.0.1 8083 127.0.0.1 8080' C-m \; \
  split-window -h \; \
  send-keys './node 127.0.0.1 8084 127.0.0.1 8080' C-m \; \
  select-pane -t 1 \; \
  send-keys './node 127.0.0.1 8085 127.0.0.1 8080' C-m \; \
  split-window -h \; \
  send-keys './node 127.0.0.1 8086 127.0.0.1 8080' C-m \; \
  select-pane -t 0 \; \
  split-window -h \; \
  send-keys './node 127.0.0.1 8087 127.0.0.1 8080' C-m \; \
  split-window -h \; \
  send-keys './node 127.0.0.1 8088 127.0.0.1 8080' C-m \; \
  select-pane -t 0 \; \
  split-window -h \; \
  send-keys './node 127.0.0.1 8089 127.0.0.1 8080' C-m \;
