PKG=svcbundle
ARCH=i386
PKGFILE=$PKG-$ARCH.pkg
DIR="./"
VERSION="1.85"

#replace version
#sed -i "s;VERSION_NUMBER;$VERSION;g" pkginfo

pkgmk -o -d $DIR -a $ARCH
touch $PKGFILE
pkgtrans -s $DIR $PKGFILE $PKG 

#file fixing
mv $PKGFILE $PKG-$VERSION-$ARCH.pkg
chmod 777 $PKG-$VERSION-$ARCH.pkg
