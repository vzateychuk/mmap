package ru.vez.kvstore;

// file: FlatKVStore.java

import com.google.flatbuffers.FlatBufferBuilder;
import java.io.File;
import java.io.IOException;
import java.io.RandomAccessFile;
import java.nio.MappedByteBuffer;
import java.nio.channels.FileChannel;
import java.util.concurrent.atomic.AtomicInteger;


public class FlatKVStore {
    private static final int MAX_RECORDS = 1024;
    private static final int MAX_BUFFER_SIZE = 1 * 1024 * 1024; // 1MB
    private static final int SLOT_SIZE = MAX_BUFFER_SIZE / MAX_RECORDS;

    private final MappedByteBuffer mmap;
    private final AtomicInteger writeIndex = new AtomicInteger(0);

    public FlatKVStore(String path) throws IOException {
        File file = new File(path);
        if (!file.exists()) {
            try (RandomAccessFile raf = new RandomAccessFile(file, "rw")) {
                raf.setLength(MAX_BUFFER_SIZE);
            }
        }
        try (RandomAccessFile raf = new RandomAccessFile(file, "rw")) {
            mmap = raf.getChannel().map(FileChannel.MapMode.READ_WRITE, 0, MAX_BUFFER_SIZE);
        }
    }

    public void put(String sender, String receiver, long amount) {
        FlatBufferBuilder builder = new FlatBufferBuilder(128);
        int senderOffset = builder.createString(sender);
        int receiverOffset = builder.createString(receiver);
        int payment = PaymentRecord.createPaymentRecord(builder, senderOffset, receiverOffset, amount);
        builder.finish(payment);

        int index = writeIndex.getAndIncrement() % MAX_RECORDS;
        int offset = index * SLOT_SIZE;

        byte[] bytes = builder.sizedByteArray();
        mmap.position(offset);
        mmap.putInt(bytes.length); // store length first
        mmap.put(bytes);
    }

    public PaymentRecord get(int index) {
        int offset = index * SLOT_SIZE;
        mmap.position(offset);
        int length = mmap.getInt();
        byte[] buf = new byte[length];
        mmap.get(buf);
        return PaymentRecord.getRootAsPaymentRecord(java.nio.ByteBuffer.wrap(buf));
    }

    public int getWriteIndex() {
        return writeIndex.get();
    }

    public static void main(String[] args) throws IOException {
        FlatKVStore store = new FlatKVStore("/tmp/payment_store.dat");

        store.put("alice", "bob", 1000L);
        store.put("carol", "dave", 2500L);

        for (int i = 0; i < store.getWriteIndex(); i++) {
            PaymentRecord record = store.get(i);
            System.out.printf("%d: %s -> %s : %d\n", i, record.sender(), record.receiver(), record.amount());
        }
    }
}
