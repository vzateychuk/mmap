# Memory-mapping golang impl

Here’s a minimal, working memory-mapped file example in Go, using golang.org/x/exp/mmap for reading or syscall for read/write.
Instead of reading and writing with traditional I/O calls, the file is accessed as if it were part of the program's memory.

## ⚙️ How It Works
Conceptual Process:

The OS maps a file into virtual memory.

The program accesses file contents like normal memory (ptr[i]), without I/O system calls.

The OS handles lazy loading (only loading parts when accessed) and flushing changes back to disk.

✅ Efficient: Reduces copying, page-level access.
🔄 Bidirectional: You can read/write memory, and it reflects to the file (if set up to do so).
🔁 Shareable: MMFs can be shared across processes for IPC (Inter-Process Communication).

## Examples:

### mmap_read_example

This example writes "Hello from mmap!" to file example.dat and reads it back via mapped memory.

```bash
go run mmap_read_example.go
```

### mmap_channel_example (write/read)

This example created a simple shared memory channel using mmap for IPC in Go:

** mmap_channel_write_process.go **- writes messages to a shared memory file.
** mmap_channel_read_process.go ** - polls the memory for new messages.

Run in two terminals:
```bash
go run mmap_channel_write_process.go
go run mmap_channel_read_process.go
```