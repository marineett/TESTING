#!/usr/bin/env python3
import argparse
import os
from pathlib import Path
from typing import Optional

import pandas as pd
import matplotlib.pyplot as plt


def read_cpu_csv(path: Path) -> Optional[pd.DataFrame]:
	if not path.exists():
		return None
	df = pd.read_csv(
		path,
		header=None,
		names=["name", "timestamp", "value"],
		quotechar='"',
	)
	df["timestamp"] = pd.to_datetime(pd.to_numeric(df["timestamp"], errors="coerce"), unit="s", utc=True)
	df["value"] = pd.to_numeric(df["value"], errors="coerce")
	df = df.dropna(subset=["timestamp", "value"])
	return df


def read_memory_csv(path: Path) -> Optional[pd.DataFrame]:
	if not path.exists():
		return None
	df = pd.read_csv(
		path,
		header=None,
		names=["name", "timestamp", "value"],
		quotechar='"',
	)
	df["timestamp"] = pd.to_datetime(pd.to_numeric(df["timestamp"], errors="coerce"), unit="s", utc=True)
	df["value"] = pd.to_numeric(df["value"], errors="coerce")
	df = df.dropna(subset=["timestamp", "value"])
	return df


def read_latency_csv(path: Path) -> Optional[pd.DataFrame]:
	if not path.exists():
		return None
	# Try header, fallback to no header
	try:
		df = pd.read_csv(path)
	except Exception:
		df = pd.read_csv(path, header=None)
	# Normalize schemas:
	cols = [str(c).strip().lower() for c in df.columns]
	if len(cols) == 3 and set(cols) >= {"chat_id", "message_index", "latency_ms"}:
		df.columns = ["chat_id", "message_index", "latency_ms"]
		df["latency_ms"] = pd.to_numeric(df.get("latency_ms"), errors="coerce")
		df["message_index"] = pd.to_numeric(df.get("message_index"), errors="coerce")
		df["chat_id"] = pd.to_numeric(df.get("chat_id"), errors="coerce")
		df = df.dropna(subset=["latency_ms"])
		return df
	if len(cols) == 2:
		# Interpret as (batch_size, latency_ms) possibly with or without header
		if "batch_size" in cols and "latency_ms" in cols:
			df = df.rename(columns={df.columns[0]: "batch_size", df.columns[1]: "latency_ms"})
		else:
			df.columns = ["batch_size", "latency_ms"]
		df["batch_size"] = pd.to_numeric(df.get("batch_size"), errors="coerce")
		df["latency_ms"] = pd.to_numeric(df.get("latency_ms"), errors="coerce")
		df = df.dropna(subset=["batch_size", "latency_ms"])
		return df
	# Fallback: just try to coerce latency_ms if present
	if "latency_ms" in df.columns:
		df["latency_ms"] = pd.to_numeric(df.get("latency_ms"), errors="coerce")
		df = df.dropna(subset=["latency_ms"])
	return df


def plot_cpu(df: pd.DataFrame, out_path: Path, title_suffix: str):
	plt.style.use("seaborn-v0_8")
	fig, ax = plt.subplots(figsize=(10, 5))
	for name, g in df.groupby("name"):
		ax.plot(g["timestamp"], g["value"], label=str(name))
	ax.set_title(f"CPU usage (cores){title_suffix}")
	ax.set_xlabel("Time")
	ax.set_ylabel("CPU (cores)")
	ax.grid(True, alpha=0.3)
	if df["name"].nunique() > 1:
		ax.legend()
	fig.autofmt_xdate()
	fig.tight_layout()
	fig.savefig(out_path, dpi=150)
	plt.close(fig)


def plot_memory(df: pd.DataFrame, out_path: Path, title_suffix: str):
	plt.style.use("seaborn-v0_8")
	fig, ax = plt.subplots(figsize=(10, 5))
	for name, g in df.groupby("name"):
		ax.plot(g["timestamp"], g["value"] / (1024**2), label=str(name))
	ax.set_title(f"Memory working set (MiB){title_suffix}")
	ax.set_xlabel("Time")
	ax.set_ylabel("Memory (MiB)")
	ax.grid(True, alpha=0.3)
	if df["name"].nunique() > 1:
		ax.legend()
	fig.autofmt_xdate()
	fig.tight_layout()
	fig.savefig(out_path, dpi=150)
	plt.close(fig)


def plot_latency(df: pd.DataFrame, out_hist: Path, out_cdf: Path, title_suffix: str):
	plt.style.use("seaborn-v0_8")
	# Histogram
	fig, ax = plt.subplots(figsize=(10, 5))
	ax.hist(df["latency_ms"], bins=50, color="#4C78A8", edgecolor="white")
	ax.set_title(f"Request latency distribution (ms){title_suffix}")
	ax.set_xlabel("Latency (ms)")
	ax.set_ylabel("Count")
	ax.grid(True, alpha=0.3)
	fig.tight_layout()
	fig.savefig(out_hist, dpi=150)
	plt.close(fig)

	lat = df["latency_ms"].sort_values().to_numpy()
	y = (pd.Series(range(1, len(lat) + 1)) / len(lat)).to_numpy()
	fig, ax = plt.subplots(figsize=(10, 5))
	ax.plot(lat, y, color="#F58518")
	ax.set_title(f"Request latency CDF{title_suffix}")
	ax.set_xlabel("Latency (ms)")
	ax.set_ylabel("Fraction <= x")
	ax.grid(True, alpha=0.3)
	fig.tight_layout()
	fig.savefig(out_cdf, dpi=150)
	plt.close(fig)

def plot_latency_vs_batch(df: pd.DataFrame, out_path: Path, title_suffix: str):
	# Expects columns: batch_size, latency_ms
	if not {"batch_size", "latency_ms"}.issubset(set(df.columns)):
		return
	g = df.groupby("batch_size", as_index=False)["latency_ms"].mean().sort_values("batch_size")
	plt.style.use("seaborn-v0_8")
	fig, ax = plt.subplots(figsize=(10, 5))
	ax.plot(g["batch_size"], g["latency_ms"], marker="o")
	ax.set_title(f"Average latency vs batch size{title_suffix}")
	ax.set_xlabel("Batch size")
	ax.set_ylabel("Average latency (ms)")
	ax.grid(True, alpha=0.3)
	fig.tight_layout()
	fig.savefig(out_path, dpi=150)
	plt.close(fig)


def main():
	parser = argparse.ArgumentParser(description="Plot metrics (CPU, memory, latency) from /metrics/<label> CSVs")
	parser.add_argument("--metrics-dir", default="metrics", help="Path to metrics root directory (default: metrics)")
	parser.add_argument("--labels", nargs="*", default=None, help="Specific labels to plot (e.g., create_batch one_big)")
	args = parser.parse_args()

	root = Path(args.metrics_dir)
	if not root.exists():
		raise SystemExit(f"Metrics directory not found: {root}")

	# Determine labels (subdirectories)
	if args.labels:
		labels = args.labels
	else:
		labels = sorted([p.name for p in root.iterdir() if p.is_dir()])

	# Collect latency stats per label for latency vs batch-size plot
	latency_summary = []

	for label in labels:
		label_dir = root / label
		if not label_dir.is_dir():
			continue
		print(f"Processing label: {label}")

		title_suffix = f" — {label}"

		cpu_csv = label_dir / "cpu.csv"
		mem_csv = label_dir / "memory.csv"
		lat_csv = label_dir / "latency_ms.csv"

		# CPU
		cpu_df = read_cpu_csv(cpu_csv)
		if cpu_df is not None and not cpu_df.empty:
			out = label_dir / "cpu.png"
			plot_cpu(cpu_df, out, title_suffix)
			print(f"Saved: {out}")
		else:
			print(f"Skip CPU (no data): {cpu_csv}")

		# Memory
		mem_df = read_memory_csv(mem_csv)
		if mem_df is not None and not mem_df.empty:
			out = label_dir / "memory.png"
			plot_memory(mem_df, out, title_suffix)
			print(f"Saved: {out}")
		else:
			print(f"Skip memory (no data): {mem_csv}")

		# Latency
		lat_df = read_latency_csv(lat_csv)
		if lat_df is not None and not lat_df.empty:
			# If file is in (batch_size, latency_ms) format – build requested plot
			if {"batch_size", "latency_ms"}.issubset(set(lat_df.columns)) and lat_df.shape[1] == 2:
				out_vs = label_dir / "latency_vs_batch.png"
				plot_latency_vs_batch(lat_df, out_vs, title_suffix)
				print(f"Saved: {out_vs}")
				# Add summaries per batch_size
				for bs, g in lat_df.groupby("batch_size"):
					latency_summary.append({
						"label": f"{label}_bs{int(bs)}",
						"batch_size": int(bs),
						"p50": float(g["latency_ms"].quantile(0.50)),
						"p95": float(g["latency_ms"].quantile(0.95)),
						"p99": float(g["latency_ms"].quantile(0.99)),
						"avg": float(g["latency_ms"].mean()),
					})
			else:
				# Otherwise, default per-request latency plots
				out_hist = label_dir / "latency_hist.png"
				out_cdf = label_dir / "latency_cdf.png"
				plot_latency(lat_df, out_hist, out_cdf, title_suffix)
				print(f"Saved: {out_hist}, {out_cdf}")
				# Summaries for cross-label plot:
				batch_size = int(lat_df["chat_id"].nunique()) if "chat_id" in lat_df.columns else None
				p50 = float(lat_df["latency_ms"].quantile(0.50))
				p95 = float(lat_df["latency_ms"].quantile(0.95))
				p99 = float(lat_df["latency_ms"].quantile(0.99))
				avg = float(lat_df["latency_ms"].mean())
				latency_summary.append({
					"label": label,
					"batch_size": batch_size,
					"p50": p50,
					"p95": p95,
					"p99": p99,
					"avg": avg,
				})
		else:
			print(f"Skip latency (no data): {lat_csv}")

	# Plot latency vs batch size (cross-label)
	if latency_summary:
		df = pd.DataFrame(latency_summary)
		# Fallback if batch_size missing
		if "batch_size" not in df or df["batch_size"].isna().all():
			df["batch_size"] = range(1, len(df) + 1)
		df = df.sort_values(by="batch_size")

		plt.style.use("seaborn-v0_8")
		fig, ax = plt.subplots(figsize=(10, 6))
		ax.plot(df["batch_size"], df["p50"], marker="o", label="p50")
		ax.plot(df["batch_size"], df["p95"], marker="o", label="p95")
		ax.plot(df["batch_size"], df["p99"], marker="o", label="p99")
		ax.plot(df["batch_size"], df["avg"], marker="o", label="avg")
		for _, row in df.iterrows():
			ax.annotate(str(row["label"]), (row["batch_size"], row["p95"]), textcoords="offset points", xytext=(0,8), ha="center", fontsize=8)
		ax.set_title("Latency vs Batch Size (per label)")
		ax.set_xlabel("Batch size (unique chats per run)")
		ax.set_ylabel("Latency (ms)")
		ax.grid(True, alpha=0.3)
		ax.legend()
		out = root / "latency_vs_batch.png"
		fig.tight_layout()
		fig.savefig(out, dpi=150)
		plt.close(fig)
		print(f"Saved: {out}")


if __name__ == "__main__":
	main()


