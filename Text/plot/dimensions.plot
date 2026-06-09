# Настройка формата вывода (окно на Mac)
set terminal qt font "Consolas,22" size 1000,600

# Если хотите сохранить в файл, раскомментируйте следующие 2 строки:
# set terminal pngcairo size 800,600 enhanced font "Arial,12"
# set output 'performance_plot.png'

set xlabel "Размер области"
set ylabel "Время выполнения (mcs)"

# Настройка сетки
set grid
set key left top  # Расположение легенды

# Определяем масштаб по оси X (чтобы метки были красивыми)
set xtics ( "1*PI" 1, "2*PI" 2, "4*PI" 4 )
set xrange [0.5:4.5]

# Рисуем графики
# Используем "-" для ввода данных прямо в скрипте
plot "-" with linespoints lw 2 ps 1.5 title "1D", \
     "-" with linespoints lw 2 ps 1.5 title "2D", \
     "-" with linespoints lw 2 ps 1.5 title "3D"

# Данные для 1D
1 14
2 37
4 109
e

# Данные для 2D
1 79
2 147
4 336
e

# Данные для 3D
1 643
2 755
4 1595
e

# Пауза, чтобы окно не закрылось сразу (при запуске из терминала)
pause -1
