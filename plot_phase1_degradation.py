#!/usr/bin/env python3
"""
ФАЗА 1: ПОИСК ТОЧКИ ДЕГРАДАЦИИ
Постепенное увеличение нагрузки для определения точки деградации системы
"""
import pandas as pd
import matplotlib.pyplot as plt
from pathlib import Path
import numpy as np

plt.rcParams['font.family'] = 'DejaVu Sans'

# === Пути ===
BASE = Path("metrics/degradation_many_user_test")
CPU_DIR = BASE / "cpu"
MEM_DIR = BASE / "memory"
LAT_DIR = BASE / "latency_ms"
OUT = Path("analysis")
OUT.mkdir(exist_ok=True)

print("="*70)
print("ФАЗА 1: ПОИСК ТОЧКИ ДЕГРАДАЦИИ")
print("="*70)

cpu_files = sorted(CPU_DIR.glob("*.csv"))[:20]
mem_files = sorted(MEM_DIR.glob("*.csv"))[:20]
lat_files = sorted(LAT_DIR.glob("*.csv"))[:20]

if not (cpu_files and mem_files and lat_files):
    print(f"❌ Файлы не найдены!")
    print(f"   CPU: {len(cpu_files)}, MEM: {len(mem_files)}, LAT: {len(lat_files)}")
    exit(1)

print(f"✅ Найдено: CPU={len(cpu_files)}, Memory={len(mem_files)}, Latency={len(lat_files)}")

# === Агрегация по итерациям ===
def aggregate_cpu(files):
    """Среднее CPU по каждой итерации"""
    return [pd.read_csv(f, names=['name', 'timestamp', 'cpu_cores'])['cpu_cores'].mean() 
            for f in files]

def aggregate_mem(files):
    """Среднее Memory по каждой итерации (в ГБ)"""
    return [pd.read_csv(f, names=['name', 'timestamp', 'memory_bytes'])['memory_bytes'].mean() / (1024**3)
            for f in files]

def aggregate_lat(files):
    """Медиана Latency по каждой итерации"""
    return [pd.read_csv(f)['latency_ms'].median() for f in files]

cpu_data = aggregate_cpu(cpu_files)
mem_data = aggregate_mem(mem_files)
lat_data = aggregate_lat(lat_files)

# Статистика
cpu_mean = np.mean(cpu_data)
cpu_max = np.max(cpu_data)
mem_mean = np.mean(mem_data)
mem_max = np.max(mem_data)
lat_mean = np.mean(lat_data)
lat_median = np.median(lat_data)
lat_max = np.max(lat_data)

print(f"\n📊 Обработано итераций: {len(cpu_data)}")

# === ПОИСК ТОЧКИ ДЕГРАДАЦИИ ===
# Ищем точку, где задержка резко увеличивается
lat_threshold = lat_mean * 1.5  # Порог деградации: 150% от среднего
degradation_point = None

for i, lat in enumerate(lat_data):
    if lat > lat_threshold:
        degradation_point = i + 1
        break

if degradation_point:
    print(f"\n⚠️  Точка деградации найдена на итерации {degradation_point}")
    print(f"   Задержка: {lat_data[degradation_point-1]:.1f} мс (порог: {lat_threshold:.1f} мс)")
else:
    print(f"\n✅ Точка деградации не найдена (все значения ниже порога {lat_threshold:.1f} мс)")

# === ГРАФИКИ ===
fig, axes = plt.subplots(3, 1, figsize=(14, 12))
iterations = range(1, len(cpu_data) + 1)

# График 1: Задержка
axes[0].plot(iterations, lat_data, 'o-', markersize=6, linewidth=2, color='#2E86AB')
axes[0].axhline(lat_mean, color='green', linestyle='--', linewidth=2, 
                label=f'Среднее: {lat_mean:.1f} мс')
axes[0].axhline(lat_threshold, color='orange', linestyle='--', linewidth=2,
                label=f'Порог деградации: {lat_threshold:.1f} мс')
if degradation_point:
    axes[0].axvline(degradation_point, color='red', linestyle=':', linewidth=2,
                    label=f'Деградация на итерации {degradation_point}')
    axes[0].plot(degradation_point, lat_data[degradation_point-1], 'ro', markersize=10)

axes[0].set_ylabel('Задержка (мс)', fontsize=13, fontweight='bold')
axes[0].set_title('Медианная задержка при увеличении нагрузки', 
                  fontsize=14, fontweight='bold')
axes[0].legend(loc='best', fontsize=10)
axes[0].grid(True, alpha=0.3)

# График 2: CPU
axes[1].plot(iterations, cpu_data, 's-', markersize=6, linewidth=2, color='#A23B72')
axes[1].axhline(cpu_mean, color='green', linestyle='--', linewidth=2,
                label=f'Среднее: {cpu_mean:.3f} ядер')
if degradation_point:
    axes[1].axvline(degradation_point, color='red', linestyle=':', linewidth=2,
                    label=f'Деградация на итерации {degradation_point}')
axes[1].set_ylabel('Использование CPU (ядер)', fontsize=13, fontweight='bold')
axes[1].set_title('Нагрузка на процессор при увеличении нагрузки', fontsize=14, fontweight='bold')
axes[1].legend(loc='best', fontsize=10)
axes[1].grid(True, alpha=0.3)

# График 3: Память
axes[2].plot(iterations, mem_data, '^-', markersize=6, linewidth=2, color='#F18F01')
axes[2].axhline(mem_mean, color='green', linestyle='--', linewidth=2,
                label=f'Среднее: {mem_mean:.2f} ГБ')
if degradation_point:
    axes[2].axvline(degradation_point, color='red', linestyle=':', linewidth=2,
                    label=f'Деградация на итерации {degradation_point}')
axes[2].set_xlabel('Номер итерации', fontsize=13, fontweight='bold')
axes[2].set_ylabel('Потребление памяти (ГБ)', fontsize=13, fontweight='bold')
axes[2].set_title('Использование памяти при увеличении нагрузки', fontsize=14, fontweight='bold')
axes[2].legend(loc='best', fontsize=10)
axes[2].grid(True, alpha=0.3)

plt.tight_layout()
output_file = OUT / "phase1_degradation_search.png"
plt.savefig(output_file, dpi=300, bbox_inches="tight")
print(f"\n✅ График: {output_file}")

# === Дополнительный график: Тренды роста ===
fig2, ax = plt.subplots(figsize=(14, 8))

# Нормализуем данные для сравнения трендов
lat_norm = np.array(lat_data) / lat_data[0] * 100
cpu_norm = np.array(cpu_data) / cpu_data[0] * 100
mem_norm = np.array(mem_data) / mem_data[0] * 100

ax.plot(iterations, lat_norm, 'o-', markersize=6, linewidth=2, 
        color='#2E86AB', label='Задержка')
ax.plot(iterations, cpu_norm, 's-', markersize=6, linewidth=2,
        color='#A23B72', label='CPU')
ax.plot(iterations, mem_norm, '^-', markersize=6, linewidth=2,
        color='#F18F01', label='Память')

ax.axhline(100, color='gray', linestyle='--', linewidth=1, alpha=0.5)
if degradation_point:
    ax.axvline(degradation_point, color='red', linestyle=':', linewidth=2,
               label=f'Точка деградации')

ax.set_xlabel('Номер итерации', fontsize=13, fontweight='bold')
ax.set_ylabel('Относительный рост (%)', fontsize=13, fontweight='bold')
ax.set_title('Сравнительный рост метрик при увеличении нагрузки (базовая итерация = 100%)', 
             fontsize=14, fontweight='bold')
ax.legend(loc='best', fontsize=11)
ax.grid(True, alpha=0.3)

plt.tight_layout()
output_file2 = OUT / "phase1_relative_growth.png"
plt.savefig(output_file2, dpi=300, bbox_inches="tight")
print(f"✅ График: {output_file2}")

# === ОТЧЁТ ===
report_file = OUT / "phase1_report.txt"
with open(report_file, 'w', encoding='utf-8') as f:
    f.write("="*70 + "\n")
    f.write("ФАЗА 1: ПОИСК ТОЧКИ ДЕГРАДАЦИИ\n")
    f.write("="*70 + "\n\n")
    
    f.write("ЧТО ДЕЛАЛИ:\n")
    f.write("Постепенно увеличивали нагрузку на систему, запуская тесты с\n")
    f.write("возрастающим числом одновременных чатов. Цель — найти точку,\n")
    f.write("в которой производительность системы начинает деградировать.\n\n")
    
    f.write("="*70 + "\n")
    f.write("СТАТИСТИКА ПО ВСЕМ ИТЕРАЦИЯМ:\n")
    f.write("="*70 + "\n\n")
    f.write(f"CPU:\n")
    f.write(f"  Среднее:  {cpu_mean:.4f} ядер\n")
    f.write(f"  Макс:     {cpu_max:.4f} ядер\n")
    f.write(f"  Прирост:  {(cpu_max/cpu_data[0]-1)*100:.1f}% от начального значения\n\n")
    
    f.write(f"Память:\n")
    f.write(f"  Среднее:  {mem_mean:.3f} ГБ\n")
    f.write(f"  Макс:     {mem_max:.3f} ГБ\n")
    f.write(f"  Прирост:  {(mem_max/mem_data[0]-1)*100:.1f}% от начального значения\n\n")
    
    f.write(f"Задержка:\n")
    f.write(f"  Среднее:  {lat_mean:.1f} мс\n")
    f.write(f"  Медиана:  {lat_median:.1f} мс\n")
    f.write(f"  Макс:     {lat_max:.1f} мс\n")
    f.write(f"  Прирост:  {(lat_max/lat_data[0]-1)*100:.1f}% от начального значения\n\n")
    
    f.write("="*70 + "\n")
    f.write("АНАЛИЗ ТОЧКИ ДЕГРАДАЦИИ:\n")
    f.write("="*70 + "\n\n")
    
    f.write(f"Порог деградации: {lat_threshold:.1f} мс (150% от среднего)\n\n")
    
    if degradation_point:
        f.write(f"🔴 ТОЧКА ДЕГРАДАЦИИ НАЙДЕНА!\n\n")
        f.write(f"Итерация:  {degradation_point}\n")
        f.write(f"Задержка:  {lat_data[degradation_point-1]:.1f} мс\n")
        f.write(f"CPU:       {cpu_data[degradation_point-1]:.4f} ядер\n")
        f.write(f"Память:    {mem_data[degradation_point-1]:.3f} ГБ\n\n")
        
        f.write("РОСТ МЕТРИК К ТОЧКЕ ДЕГРАДАЦИИ:\n")
        f.write(f"  Задержка: +{(lat_data[degradation_point-1]/lat_data[0]-1)*100:.1f}%\n")
        f.write(f"  CPU:      +{(cpu_data[degradation_point-1]/cpu_data[0]-1)*100:.1f}%\n")
        f.write(f"  Память:   +{(mem_data[degradation_point-1]/mem_data[0]-1)*100:.1f}%\n\n")
    else:
        f.write(f"✅ ТОЧКА ДЕГРАДАЦИИ НЕ ДОСТИГНУТА\n\n")
        f.write(f"Все {len(lat_data)} итераций прошли успешно.\n")
        f.write(f"Максимальная задержка: {lat_max:.1f} мс (ниже порога {lat_threshold:.1f} мс)\n\n")
    
    f.write("="*70 + "\n")
    f.write("ВЫВОДЫ:\n")
    f.write("="*70 + "\n\n")
    
    # Анализ линейности роста
    lat_growth_rate = (lat_data[-1] / lat_data[0] - 1) / len(lat_data) * 100
    cpu_growth_rate = (cpu_data[-1] / cpu_data[0] - 1) / len(cpu_data) * 100
    mem_growth_rate = (mem_data[-1] / mem_data[0] - 1) / len(mem_data) * 100
    
    f.write("ТЕМПЫ РОСТА НА ИТЕРАЦИЮ:\n")
    f.write(f"  Задержка: {lat_growth_rate:.2f}% за итерацию\n")
    f.write(f"  CPU:      {cpu_growth_rate:.2f}% за итерацию\n")
    f.write(f"  Память:   {mem_growth_rate:.2f}% за итерацию\n\n")
    
    if degradation_point:
        recommended_load = degradation_point - 1
        f.write("РЕКОМЕНДАЦИИ:\n")
        f.write(f"✅ Рекомендуемая максимальная нагрузка: {recommended_load} итераций\n")
        f.write(f"   (на 1 итерацию меньше точки деградации)\n\n")
        
        f.write("СЛЕДУЮЩИЙ ШАГ:\n")
        f.write("📋 Фаза 2: Протестировать стабильность на найденной нагрузке\n")
        f.write(f"   (запустить {recommended_load}-ю конфигурацию 20 раз)\n")
    else:
        f.write("РЕКОМЕНДАЦИИ:\n")
        f.write(f"✅ Система справляется с максимальной тестируемой нагрузкой\n")
        f.write(f"   (все {len(lat_data)} итераций прошли успешно)\n\n")
        
        f.write("СЛЕДУЮЩИЙ ШАГ:\n")
        f.write("📋 Можно увеличить нагрузку для поиска реального предела\n")
        f.write("   или перейти к тестированию стабильности на текущем уровне\n")
    
    f.write("\n" + "="*70 + "\n")

print(f"✅ Отчёт: {report_file}\n")
print("="*70)
