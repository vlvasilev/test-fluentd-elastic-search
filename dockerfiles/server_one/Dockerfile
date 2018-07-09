FROM alpine
RUN mkdir /app
RUN mkdir /conf
ADD ./server_one .
ADD ./deploy_job.sh .
ADD ./jobtemplate.yaml ./
ADD ./kubectl /usr/local/bin/
EXPOSE 8000
CMD ["./server_one"]

