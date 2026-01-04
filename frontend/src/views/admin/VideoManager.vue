<template>
  <div class="video-manager">
    <div class="toolbar">
      <el-button type="primary" @click="handleCreate">上传视频</el-button>
    </div>

    <el-table :data="tableData" style="width: 100%" v-loading="loading">
      <el-table-column prop="opera_id" label="ID" width="80" />
      <el-table-column label="封面" width="100">
        <template #default="scope">
          <el-image :src="scope.row.avatar" style="width: 80px; height: 45px" fit="cover" />
        </template>
      </el-table-column>
      <el-table-column prop="opera_title" label="标题" />
      <el-table-column label="艺术家">
        <template #default="scope">
          <span v-for="(artist, index) in scope.row.artists" :key="artist.artist_id">
            {{ artist.name }}<span v-if="index < scope.row.artists.length - 1">, </span>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.is_hidden ? 'info' : 'success'">
            {{ scope.row.is_hidden ? '隐藏' : '公开' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" :type="scope.row.is_hidden ? 'success' : 'warning'"
            @click="handleToggleHidden(scope.row)">
            {{ scope.row.is_hidden ? '显示' : '隐藏' }}
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="total"
        layout="prev, pager, next" @current-change="fetchData" />
    </div>

    <!-- Edit/Create Dialog -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑视频' : '上传视频'" width="50%" @paste="handlePaste">
      <el-form :model="form" label-width="100px">
        <el-form-item label="标题">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="视频文件">
          <!-- Upload Logic (Simplified) -->
          <el-upload ref="videoUploadRef" class="upload-demo" action="#" :auto-upload="false" :limit="1"
            :on-change="(file: any) => handleElFileChange(file, 'video')" :on-remove="() => handleElFileRemove('video')"
            :on-exceed="handleExceed" drag>
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将视频拖到此处，或<em>点击上传</em>
            </div>
          </el-upload>
          <div v-if="videoFile" class="custom-progress-track">
            <div class="custom-progress-bar" :style="{ width: videoProgress + '%' }"></div>
          </div>
          <div v-if="isEdit" class="el-form-item__error" style="position: static; color: #909399;">如果不修改视频，请忽略此项</div>
        </el-form-item>
        <el-form-item label="封面图片">
          <el-upload ref="avatarUploadRef" class="upload-demo" action="#" :auto-upload="false" :limit="1"
            :on-change="(file: any) => handleElFileChange(file, 'avatar')"
            :on-remove="() => handleElFileRemove('avatar')" :on-exceed="handleExceed" list-type="picture" drag>
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
          <div v-if="avatarFile" class="custom-progress-track">
            <div class="custom-progress-bar" :style="{ width: avatarProgress + '%' }"></div>
          </div>
        </el-form-item>
        <el-form-item label="艺术家">
          <el-select v-model="form.artist_ids" multiple filterable allow-create default-first-option
            :reserve-keyword="false" placeholder="请选择或输入艺术家名字">
            <el-option v-for="item in artistOptions" :key="item.artist_id" :label="item.name" :value="item.artist_id" />
          </el-select>
        </el-form-item>
        <el-form-item label="隐藏">
          <el-switch v-model="form.is_hidden" />
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
import VideoManagerScript from '../../scripts/views/admin/VideoManager';
import { UploadFilled } from '@element-plus/icons-vue';

export default {
  ...VideoManagerScript,
  components: {
    UploadFilled
  }
};
</script>

<style scoped src="../../styles/views/admin/VideoManager.css"></style>
