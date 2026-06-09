set terminal qt font "Consolas,22" size 500,500

set xrange [-1:1]
set yrange [-0.05:1.1]
set samples 400

set border lw 1.5
set grid lw 1

set xlabel "x" font "Consolas,22"
set ylabel "y" font "Consolas,22"
set xtics font "Consolas,22"
set ytics font "Consolas,22"

set key bottom right font "Consolas,22" spacing 1.3 samplen 2.2

linear_hat(x) = (abs(x) <= 1 ? 1 - abs(x) : 0)
quadratic_basis(x) = (abs(x) <= 1 ? 1 - x**2 : 0)

plot linear_hat(x) lw 3 lc rgb "#1f77b4" title "Линейный базис", \
     quadratic_basis(x) lw 3 lc rgb "#d62728" title "Квадратичный базис"

pause -1
