# crazyreply

>一个疯狂发消息的应用

### ALL:
```
golang 1.15.2
git
```

#### Fedora:
```
xclip xsel要保证其中一个可以正常调用
$ dnf install gtk3-devel gdk-pixbuf2-devel glib2-devel \
    libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libxkbfile-devel
$ git clone https://github.com/xiefeihong/crazyreply.git
$ cd crazyreply
$ go build
```
    
#### Windows:
```
推荐安装msys2
$ pacman -S mingw-w64-x86_64-gtk3 mingw-w64-x86_64-toolchain base-devel glib2-devel
$ source ~/.bashrc
$ sed -i -e 's/-Wl,-luuid/-luuid/g' /mingw64/lib/pkgconfig/gdk-3.0.pc # This fixes a bug in pkgconfig
$ git clone https://github.com/xiefeihong/crazyreply.git
$ cd crazyreply
$ go build  -ldflags="-H windowsgui"
```

#### 使用方法:
```
输入预先填好的文字信息，点击开始，然后把焦点放到对应的文本框里就可以发消息了。
```
