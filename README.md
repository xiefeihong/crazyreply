# crazyreply [![GitHub release](https://img.shields.io/github/release/xiefeihong/crazyreply.svg)](https://github.com/xiefeihong/crazyreply/releases/latest)

>一个疯狂发消息的应用

### 应用说明
```
此项目是将预设好的消息通过粘贴到文本框并控制键盘的方式实现的循环发消息。
此程序依赖robotgo，gotk3；目前确定支持Windows和Linux，MacOS理论上也支持。

```

#### 使用方法:
```
输入预先填好的文字信息，点击开始，然后把焦点放到对应的文本框里就可以发消息了。
```

### 注意事项
```
此项目在运行的时候可能会因为把焦点错误指定到某个app的窗口，导致消息发送给了错误的人，所以在使用本工具之前应该谨慎关闭桌面多余的窗口。
如果使用此工具进行特殊的用途或者造成了误会和损失本人概不负责。
```

### 建需要的软件:
```
golang
git
```

#### Fedora构建:
```
xclip xsel要保证其中一个可以正常调用
$ dnf install gtk3-devel gdk-pixbuf2-devel glib2-devel \
    libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libxkbfile-devel
$ git clone https://github.com/xiefeihong/crazyreply.git
$ cd crazyreply
$ go build
```
    
#### Windows构建:
```
推荐安装msys2
$ pacman -S mingw-w64-x86_64-gtk3 mingw-w64-x86_64-toolchain base-devel glib2-devel
$ source ~/.bashrc
$ sed -i -e 's/-Wl,-luuid/-luuid/g' /mingw64/lib/pkgconfig/gdk-3.0.pc # This fixes a bug in pkgconfig
$ git clone https://github.com/xiefeihong/crazyreply.git
$ cd crazyreply
$ go build  -ldflags="-H windowsgui"
```

#### Licences

crazyreply is licenced under [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

```
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```
