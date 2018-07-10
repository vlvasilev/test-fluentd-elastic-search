FROM alpine
RUN mkdir /app & mkdir /conf
ADD ./bin/server_one/server_one .
ADD ./resourses/kubectl /usr/local/bin/
EXPOSE 8000
CMD ["./server_one"]

