## Kirin_DB

手写实现MySQL的基本功能，涵盖文件系统实现，redo实现，buffer pool实现，undo实现，事务实现，并发控制实现，锁实现，record记录实现，表扫描实现。。。

### 1.file_manager文件系统模块

### 2.log_manager日志模块实现

### 3.buffer_manager缓存模块实现

### 4.undo回滚日志实现（tx目录下）

CHECKPOINT 、START 、COMMIT 、ROLLBACK 、SETINT 、SETSTRING六种日志的实现

### 5.事务对象实现

一个事务对象包括三个部分：
</br>
BufferList:用于管理当前被pin的buffer对象
</br>
RecoverMgr:用于事务恢复或者回滚
</br>
CocurrentMgt:用于管理并发请求