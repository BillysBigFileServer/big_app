FROM golang:bookworm AS builder-android
WORKDIR /build

RUN apt update
RUN apt-get install curl unzip -y

RUN curl -LO https://dl.google.com/android/repository/android-ndk-r27b-linux.zip && unzip android-ndk-r27b-linux.zip
ENV ANDROID_NDK_HOME=/build/android-ndk-r27b/

RUN apt-get install golang gcc gcc-mingw-w64 -y
RUN go install fyne.io/fyne/v2/cmd/fyne@latest

COPY go.mod go.sum main.go Icon.png .
COPY helper ./helper
COPY pages ./pages
COPY preferences ./preferences
COPY resources ./resources
COPY state ./state

RUN fyne package --target android -appID io.bbfs.app -icon Icon.png --release



FROM golang:alpine AS builder-linux
WORKDIR /build

RUN apk add go gcc libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev linux-headers mesa-dev xz
RUN go install fyne.io/fyne/v2/cmd/fyne@latest

COPY go.mod go.sum main.go Icon.png .
COPY helper ./helper
COPY pages ./pages
COPY preferences ./preferences
COPY resources ./resources
COPY state ./state

RUN fyne package --target linux -appID io.bbfs.app -icon Icon.png --release  && mv app.tar.xz app_linux.tar.xz


FROM scratch
WORKDIR /out
COPY --from=builder-android /build/app.apk .
COPY --from=builder-linux /build/app_linux.tar.xz .
