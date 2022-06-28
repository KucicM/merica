from cProfile import label
from dataclasses import dataclass
from typing import Tuple
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
from mpl_toolkits.mplot3d import Axes3D


@dataclass(init=False)
class Record:
    bench_type: str
    impl: str
    reader_count: int
    writer_count: int
    cpu_count: int
    number_of_iterations: int
    time_per_ops: float
    time_per_ops_unit: str
    memory_per_ops: int
    memory_per_ops_unit: str
    memory_allocs_per_ops: int

    def __init__(self, line: str):
        bench_name, rest = line.split("/", 1)
        self.bench_type = parseBenchType(bench_name)
        self.impl = parseImpl(bench_name, self.bench_type)

        trial, iter_count, time_ops, mem_ops, mem_allocs = rest.split("\t")
        self.writer_count, self.reader_count = parseWRCounts(trial)
        self.cpu_count = parseCPUCount(trial)
        self.number_of_iterations = parseIterCount(iter_count)
        self.time_per_ops, self.time_per_ops_unit = parseOpsTimeAndUnit(time_ops)
        self.memory_per_ops, self.memory_per_ops_unit = parseMemOpsAndUnit(mem_ops)
        self.memory_allocs_per_ops = parseMemAllocs(mem_allocs)


def parseBenchType(name: str) -> str:
    if name.endswith("Parallel"):
        return "Parallel"
    elif name.endswith("Sequential"):
        return "Sequential"
    else:
        raise RuntimeError(f"Check type in name {name}")


def parseImpl(name: str, bench_type: str) -> str:
    return name.replace("Benchmark", "").replace(bench_type, "")


def parseWRCounts(trial: str) -> Tuple[int, int]:
    # Writers:1_Readers:6-12 -> (1, 6)
    trial = trial.strip().split("-")[0].split("_")
    # [Writers:1,Readers:6]
    return (int(c.split(":")[1]) for c in trial)


def parseCPUCount(trial: str) -> int:
    # Writers:1_Readers:1-12 -> 12
    splt = trial.split("-")
    return 1 if len(splt) == 1 else int(splt[-1])


def parseIterCount(iter_count: str) -> int:
    return int(iter_count.strip())


def parseOpsTimeAndUnit(time_ops: str) -> Tuple[float, str]:
    # 3263 ns/op -> (3263, ns/op)
    time, unit = time_ops.strip().split(" ")
    return float(time), unit


def parseMemOpsAndUnit(mem_ops: str) -> Tuple[int, str]:
    # 999 B/op -> (999, B/op)
    mem, unit = mem_ops.strip().split(" ")
    return int(mem), unit


def parseMemAllocs(mem_allocs: str) -> int:
    # 10 allocs/op -> 10
    return int(mem_allocs.strip().split(" ")[0])


def main(data_path: str):
    records = []
    with open(data_path, "r") as f:
        while line := f.readline():
            if not line.startswith("Benchmark"):
                continue
            records.append(Record(line))

    df = pd.DataFrame(records)
    assert len(df["time_per_ops_unit"].unique()) == 1, "TODO suport diffrent time units"

    df = df.groupby(["bench_type", "impl", "reader_count", "writer_count", "cpu_count"]).median().reset_index()

    import numpy as np
    bench_types = df["bench_type"].unique()
    impls = df["impl"].unique()

    fig = plt.figure()
    for bench_type in df["bench_type"].unique():
        ax = fig.add_subplot(111, projection="3d")
        ax.set_xlabel("Number of readers")
        ax.set_ylabel("Number of writers")
        ax.set_zlabel("Time pep operation")
        for impl in df["impl"].unique():
            df_copy = df[df["bench_type"] == bench_type]
            df_copy = df_copy[df_copy["impl"] == impl]
            x = np.log2(df_copy["reader_count"])
            y = np.log2(df_copy["writer_count"])
            z = df_copy["time_per_ops"]
            ax.plot_trisurf(x, y, z, label=impl)
        break

    fig = plt.figure()
    readers = [1, 2, 8, 2**8, 2**16]
    writers = sorted(df["writer_count"].unique().tolist())
    markers = ['x', 'o', '^']
    for bench_type in bench_types:
        ax = fig.add_subplot(111)
        print(bench_type)
        for i, impl in enumerate(impls):
            marker = markers[i]
            for rds in readers:
                dfc = df[df["bench_type"] == bench_type]
                dfc = dfc[dfc["impl"] == impl]
                dfc = dfc[dfc["reader_count"] == rds]
                ax.plot(np.log2(dfc["writer_count"]), dfc["time_per_ops"], label=f"{impl} - {rds}", marker=marker)
        break
    # for bench_type in df["bench_type"].unique():
    #     ax = fig.add_subplot(111, projection="3d")
    #     ax.set_xlabel("Number of readers")
    #     ax.set_ylabel("Number of writers")
    #     ax.set_zlabel("Time pep operation")
    #     for impl in df["impl"].unique():

    # df = df.groupby(["bench_type", "impl", "reader_count", "writer_count", "cpu_count"]).median().reset_index()
    # print(df)

    # d = df[df["bench_type"] == "Parallel"]
    # impl = "WaitFreeLinkedList"
    # d = d[d["impl"] == impl]
    # x = np.log2(d["reader_count"])
    # y = np.log2(d["writer_count"])
    # z = d["time_per_ops"]


    # # ax.scatter(x, y, z, color="blue")
    # ax.plot_trisurf(x, y, z, color="blue")

    # d = df[df["bench_type"] == "Sequential"]
    # d = d[d["impl"] == impl]
    # x = np.log2(d["reader_count"])
    # y = np.log2(d["writer_count"])
    # z = d["time_per_ops"]
    # ax.plot_trisurf(x, y, z, color="red")

    plt.legend()
    plt.show()


if __name__ == "__main__":
    main("data/results_test.txt")