import pandas as pd
import matplotlib.pyplot as plt

no_index = pd.read_csv('no_index.csv', 
                      header=0,
                      names=['users', 'ms'],
                      dtype={'users': int, 'ms': int})

with_index = pd.read_csv('with_index.csv',
                        header=0,
                        names=['users', 'ms'],
                        dtype={'users': int, 'ms': int})

plt.figure(figsize=(10, 6))

plt.plot(no_index['users'], no_index['ms'], 'b-', label='Без индекса', marker='o')
plt.plot(with_index['users'], with_index['ms'], 'r-', label='С индексом', marker='s')

plt.title('Сравнение времени выполнения с индексом и без')
plt.xlabel('Количество пользователей')
plt.ylabel('Время выполнения (мс)')
plt.grid(True)
plt.legend()

plt.savefig('benchmark_comparison.png', dpi=300, bbox_inches='tight')
plt.show()