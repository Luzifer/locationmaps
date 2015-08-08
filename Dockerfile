FROM scratch

ADD ./ca-certificates.pem /etc/ssl/ca-bundle.pem
ADD ./locationmaps /locationmaps

EXPOSE 3000

ENTRYPOINT ["/locationmaps"]
CMD ["--"]
