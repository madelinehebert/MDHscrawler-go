PKG=scrawler
ARCH=i386
PKGFILE=$PKG-$ARCH.pkg
DIR="./"
VERSION=`cat pkginfo | grep VERSION | cut -d "=" -f 2`

#remove old package
rm -rf scrawler*

#
pkgmk -o -d $DIR -a $ARCH
touch $PKGFILE
pkgtrans -s $DIR $PKGFILE $PKG 

#file fixing
mv $PKGFILE $PKG-$VERSION-$ARCH.pkg
chmod 777 $PKG-$VERSION-$ARCH.pkg
