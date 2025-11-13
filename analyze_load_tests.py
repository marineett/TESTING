#!/usr/bin/env python3
"""
ФАЗА 3: Графики нагрузочного тестирования
"""
import pandas as pd
import matplotlib.pyplot as plt
from pathlib import Path
import numpy as np

plt.rcParams['font.family'] = 'DejaVu Sans'

BASE = Path("metrics/degradation_many_user_test")
CPU_DIR = BASE / "cpu"
MEM_DIR = BASE / "memory"
LAT_DIR = BASE / "latency_ms"
OUT = Path("analysis")
OUT.mkdir(exist_ok=True)

print("Чтение данных...")

cpu_files = sorted(CPU_DIR.glob("*.csv"))[:20]
mem_files = sorted(MEM_DIR.glob("*.csv"))[:20]
lat_files = sorted(LAT_DIR.glob("*.csv"))[:20]

# === АГРЕГАЦИЯ ПО ИТЕРАЦИЯМ (не по временным точкам) ===
def aggregate_cpu(files):
    """Среднее CPU по каждой итерации"""
    return [pd.read_csv(f, names=['name', 'timestamp', 'cpu_cores'])['cpu_cores'].mean() 
            for f in files]

def aggregate_mem(files):
    """Среднее Memory по каждой итерации"""
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
cpu_std = np.std(cpu_data)
mem_mean = np.mean(mem_data)
mem_std = np.std(mem_data)
lat_mean = np.mean(lat_data)
lat_median = np.median(lat_data)
lat_std = np.std(lat_data)

print(f"Загружено итераций: {len(cpu_data)}")

# === Графики ===
fig, axes = plt.subplots(3, 1, figsize=(14, 12))
iterations = range(1, len(cpu_data) + 1)

# График 1: Задержка
axes[0].plot(iterations, lat_data, 'o-', markersize=6, linewidth=2, color='#2E86AB')
axes[0].axhline(lat_mean, color='red', linestyle='--', linewidth=2, 
                label=f'Среднее: {lat_mean:.1f} мс')
axes[0].set_ylabel('Задержка (мс)', fontsize=13, fontweight='bold')
axes[0].set_title('Медианная задержка по итерациям (80 чатов)', 
                  fontsize=14, fontweight='bold')
axes[0].legend(loc='best', fontsize=11)
axes[0].grid(True, alpha=0.3)

# График 2: CPU
axes[1].plot(iterations, cpu_data, 's-', markersize=6, linewidth=2, color='#A23B72')
axes[1].axhline(cpu_mean, color='red', linestyle='--', linewidth=2,
                label=f'Среднее: {cpu_mean:.3f} ядер')
axes[1].set_ylabel('Использование CPU (ядер)', fontsize=13, fontweight='bold')
axes[1].set_title('Средняя нагрузка на процессор', fontsize=14, fontweight='bold')
axes[1].legend(loc='best', fontsize=11)
axes[1].grid(True, alpha=0.3)

# График 3: Память
axes[2].plot(iterations, mem_data, '^-', markersize=6, linewidth=2, color='#F18F01')
axes[2].axhline(mem_mean, color='red', linestyle='--', linewidth=2,
                label=f'Среднее: {mem_mean:.2f} ГБ')
axes[2].set_xlabel('Номер итерации', fontsize=13, fontweight='bold')
axes[2].set_ylabel('Потребление памяти (ГБ)', fontsize=13, fontweight='bold')
axes[2].set_title('Среднее использование памяти', fontsize=14, fontweight='bold')
axes[2].legend(loc='best', fontsize=11)
axes[2].grid(True, alpha=0.3)

plt.tight_layout()
output_file = OUT / "phase3_load_test_80chats.png"
plt.savefig(output_file, dpi=300, bbox_inches="tight")
print(f"\n✅ График: {output_file}")

# === ОТЧЁТ ===
report_file = OUT / "phase3_report.txt"
with open(report_file, 'w', encoding='utf-8') as f:
    f.write("="*70 + "\n")
    f.write("ФАЗА 3: ТЕСТИРОВАНИЕ НА НАГРУЗКЕ 80 ЧАТОВ\n")
    f.write("="*70 + "\n\n")
    
    f.write("📊 СТАТИСТИКА:\n\n")
    f.write(f"CPU:       {cpu_mean:.4f} ядер (σ={cpu_std:.4f}, вариация={cpu_std/cpu_mean*100:.1f}%)\n")
    f.write(f"Память:    {mem_mean:.3f} ГБ (σ={mem_std:.3f}, вариация={mem_std/mem_mean*100:.1f}%)\n")
    f.write(f"Задержка:  {lat_mean:.1f} мс (медиана={lat_median:.1f}, σ={lat_std:.1f}, вариация={lat_std/lat_mean*100:.1f}%)\n\n")
    
    throughput = 1000.0 / lat_mean
    f.write(f"Пропускная способность: {throughput:.2f} req/s\n\n")
    
    f.write("="*70 + "\n")
    f.write("ВЫВОДЫ:\n\n")
    
    stable = (cpu_std/cpu_mean < 0.1 and mem_std/mem_mean < 0.1 and lat_std/lat_mean < 0.15)
    
    if stable:
        f.write("✅ Система СТАБИЛЬНА на нагрузке 80 чатов\n")
        f.write("   Все метрики имеют низкую вариацию (<10-15%)\n")
    else:
        f.write("⚠️  Система показывает флуктуации\n")
        f.write("   Некоторые метрики нестабильны\n")

print(f"✅ Отчёт: {report_file}")