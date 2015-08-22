FROM scratch

ADD https://rootcastore.hub.luzifer.io/v1/store/latest /etc/ssl/ca-bundle.pem
ADD ./locationmaps /locationmaps

EXPOSE 3000

ENTRYPOINT ["/locationmaps"]
CMD ["--"]
