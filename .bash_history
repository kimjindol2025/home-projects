claude
exit
# ~/.bashrc에 alias 추가
echo 'alias claude="proot --bind=$PREFIX/tmp:/tmp $(which claude)"' >> ~/.bashrc
source ~/.bashrc
# 확인
grep claude ~/.bashrc
exit
# proot 설치
pkg install proot -y
# /tmp 마운트해서 실행
proot --bind=$PREFIX/tmp:/tmp claude
exit
ls
tail -f /tmp/claude-10603/-data-data-com-termux-files-
ks googs
exit
kc f
ks gogs
exit
