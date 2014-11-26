FROM google/golang
WORKDIR /gopath
ADD . /gopath
CMD ["bash"]
