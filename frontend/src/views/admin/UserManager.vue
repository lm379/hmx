<template>
  <div class="user-manager">
    <el-table :data="tableData" style="width: 100%" v-loading="loading">
      <el-table-column prop="user_id" label="ID" width="80" />
      <el-table-column label="头像" width="80">
        <template #default="scope">
           <el-avatar :src="scope.row.icon || `https://ui-avatars.com/api/?name=${scope.row.username}&background=random`" />
        </template>
      </el-table-column>
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="phone" label="手机号" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column prop="role" label="角色">
         <template #default="scope">
             <el-tag :type="scope.row.role === 'Administrator' ? 'danger' : 'primary'">{{ scope.row.role }}</el-tag>
         </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
            <el-popconfirm 
                :title="scope.row.role === 'User' ? '设为管理员?' : '设为普通用户?'"
                @confirm="toggleRole(scope.row)"
            >
                <template #reference>
                    <el-button size="small" :type="scope.row.role === 'User' ? 'warning' : 'info'">
                        {{ scope.row.role === 'User' ? '设为管理员' : '取消管理员' }}
                    </el-button>
                </template>
            </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchData"
      />
    </div>
  </div>
</template>

<script lang="ts">
import UserManagerScript from '../../scripts/views/admin/UserManager';
export default UserManagerScript;
</script>

<style scoped src="../../styles/views/admin/UserManager.css"></style>
