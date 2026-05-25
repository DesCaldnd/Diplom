# Настройка терминала
set terminal qt font "Consolas,16" size 1000,600
set title "Время выполнения (mcs)" font "Consolas,16"
set ylabel "Микросекунды (mcs)"
set grid y

# --- ИСПРАВЛЕНИЯ ЗДЕСЬ ---
set yrange [0:*]       # Начать отсчет Y с нуля
set xrange [-0.5:1.5]  # Центрирование (для 2 групп данных)
# -------------------------

# Настройка стиля гистограмм
set style data histograms
set style histogram cluster gap 1
set style fill solid 0.7 border -1
set boxwidth 0.8

# Определяем блок данных заранее
$MyData << EOD
Label         Linear     Quadratic
Parallel      557362     511620
Sequential    1920580    1699398
EOD

# Рисуем
plot $MyData using 2:xtic(1) title "Linear", \
     '' using 3 title "Quadratic"

pause -1
