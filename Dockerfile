FROM scratch

ADD ./ca-certificates.pem /etc/ssl/ca-bundle.pem
ADD ./locationmaps /locationmaps

ENTRYPOINT ["/locationmaps"]
CMD ["--"]
