<template>
  <div class="artist-manager">
    <div class="toolbar">
      <el-button type="primary" @click="handleCreate">添加艺术家</el-button>
    </div>

    <el-table :data="tableData" style="width: 100%" v-loading="loading">
      <el-table-column prop="artist_id" label="ID" width="80" />
      <el-table-column label="头像" width="100">
        <template #default="scope">
          <el-image :src="scope.row.avatar" style="width: 50px; height: 50px; border-radius: 50%" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="name" label="姓名" />
      <el-table-column prop="bio" label="简介" show-overflow-tooltip />
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
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

    <!-- Dialog -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑艺术家' : '添加艺术家'" width="40%" @paste="handlePaste">
      <el-form :model="form" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.bio" type="textarea" />
        </el-form-item>
        <el-form-item label="头像">
            <el-upload
               ref="avatarUploadRef"
               class="upload-demo"
               action="#"
               :auto-upload="false"
               :limit="1"
               :on-change="handleElFileChange"
               :on-remove="handleElFileRemove"
               :on-exceed="handleExceed"
               list-type="picture"
               drag
             >
               <el-icon class="el-icon--upload"><upload-filled /></el-icon>
               <div class="el-upload__text">
                 将文件拖到此处，或<em>点击上传</em>
               </div>
               <template #tip>
                 <div class="el-upload__tip">
                   支持拖拽图片，或在对话框内粘贴剪贴板图片
                 </div>
               </template>
             </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import ArtistManagerScript from '../../scripts/views/admin/ArtistManager';
import { UploadFilled } from '@element-plus/icons-vue';

export default {
  ...ArtistManagerScript,
  components: {
    UploadFilled
  }
};
</script>

<style scoped src="../../styles/views/admin/ArtistManager.css"></style>
