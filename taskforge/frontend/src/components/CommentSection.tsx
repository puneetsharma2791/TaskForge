import { useState } from 'react';
import { IconSend, IconEdit, IconTrash, IconX, IconCheck } from '@tabler/icons-react';
import { commentsService } from '@/services/comments.service';
import { useApi } from '@/hooks/useApi';
import { useMutation } from '@/hooks/useMutation';
import type { Comment } from '@/types';

const MAX_COMMENT_LENGTH = 2500;

export default function CommentSection({ taskId }: { taskId: string }) {
  const [newComment, setNewComment] = useState('');
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editContent, setEditContent] = useState('');
  const [deleteConfirmId, setDeleteConfirmId] = useState<string | null>(null);

  const { data: comments, loading, refetch } = useApi(() =>
    commentsService.list(taskId)
  );

  const { mutate: addComment, loading: adding } = useMutation(
    (content: string) => commentsService.add(taskId, content),
    { onSuccess: () => { setNewComment(''); refetch(); } }
  );

  const { mutate: editComment } = useMutation(
    ({ commentId, content }: { commentId: string; content: string }) =>
      commentsService.edit(taskId, commentId, content),
    { onSuccess: () => { setEditingId(null); setEditContent(''); refetch(); } }
  );

  const { mutate: deleteComment } = useMutation(
    (commentId: string) => commentsService.delete(taskId, commentId),
    { onSuccess: () => { setDeleteConfirmId(null); refetch(); } }
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = newComment.trim();
    if (!trimmed) return;
    addComment(trimmed);
  };

  const handleEdit = (comment: Comment) => {
    setEditingId(comment.id);
    setEditContent(comment.content);
    setDeleteConfirmId(null);
  };

  const handleSaveEdit = (commentId: string) => {
    const trimmed = editContent.trim();
    if (!trimmed) return;
    editComment({ commentId, content: trimmed });
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setEditContent('');
  };

  const handleDeleteClick = (commentId: string) => {
    setDeleteConfirmId(commentId);
    setEditingId(null);
  };

  const handleConfirmDelete = (commentId: string) => {
    deleteComment(commentId);
  };

  if (loading) {
    return (
      <div className="flex justify-center p-4">
        <span className="loading loading-spinner loading-sm"></span>
      </div>
    );
  }

  return (
    <div className="mt-6">
      <h3 className="text-lg font-semibold mb-4">
        Comments {comments && comments.length > 0 && `(${comments.length})`}
      </h3>

      {/* Add comment form */}
      <form onSubmit={handleSubmit} className="mb-4">
        <div className="flex gap-2">
          <div className="flex-1">
            <textarea
              className="textarea textarea-bordered w-full"
              placeholder="Add a comment..."
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              maxLength={MAX_COMMENT_LENGTH}
              rows={2}
            />
            {newComment.length > 0 && (
              <div className="text-xs text-base-content/50 text-right mt-1">
                {newComment.length}/{MAX_COMMENT_LENGTH}
              </div>
            )}
          </div>
          <button
            type="submit"
            className="btn btn-primary btn-sm self-start mt-1"
            disabled={!newComment.trim() || adding}
          >
            {adding ? (
              <span className="loading loading-spinner loading-xs"></span>
            ) : (
              <IconSend size={16} />
            )}
          </button>
        </div>
      </form>

      {/* Comments list */}
      <div className="space-y-3">
        {(!comments || comments.length === 0) && (
          <p className="text-sm text-base-content/50">No comments yet.</p>
        )}
        {comments?.map((comment) => (
          <div key={comment.id} className="card bg-base-200 p-3">
            <div className="flex items-start justify-between">
              <div className="text-xs text-base-content/60 mb-1">
                <span className="font-medium">{comment.author_id}</span>
                <span className="mx-1">&middot;</span>
                <span>{new Date(comment.created_at).toLocaleString()}</span>
                {comment.updated_at !== comment.created_at && (
                  <span className="italic ml-1">(edited)</span>
                )}
              </div>
              {editingId !== comment.id && deleteConfirmId !== comment.id && (
                <div className="flex gap-1">
                  <button
                    className="btn btn-ghost btn-xs"
                    onClick={() => handleEdit(comment)}
                    title="Edit"
                  >
                    <IconEdit size={14} />
                  </button>
                  <button
                    className="btn btn-ghost btn-xs text-error"
                    onClick={() => handleDeleteClick(comment.id)}
                    title="Delete"
                  >
                    <IconTrash size={14} />
                  </button>
                </div>
              )}
            </div>

            {editingId === comment.id ? (
              <div className="mt-1">
                <textarea
                  className="textarea textarea-bordered w-full text-sm"
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  maxLength={MAX_COMMENT_LENGTH}
                  rows={2}
                />
                <div className="flex gap-1 mt-1">
                  <button
                    className="btn btn-primary btn-xs"
                    onClick={() => handleSaveEdit(comment.id)}
                    disabled={!editContent.trim()}
                  >
                    <IconCheck size={14} /> Save
                  </button>
                  <button
                    className="btn btn-ghost btn-xs"
                    onClick={handleCancelEdit}
                  >
                    <IconX size={14} /> Cancel
                  </button>
                </div>
              </div>
            ) : deleteConfirmId === comment.id ? (
              <div className="mt-1">
                <p className="text-sm text-warning mb-2">Delete this comment?</p>
                <div className="flex gap-1">
                  <button
                    className="btn btn-error btn-xs"
                    onClick={() => handleConfirmDelete(comment.id)}
                  >
                    Confirm Delete
                  </button>
                  <button
                    className="btn btn-ghost btn-xs"
                    onClick={() => setDeleteConfirmId(null)}
                  >
                    Cancel
                  </button>
                </div>
              </div>
            ) : (
              <p className="text-sm whitespace-pre-wrap">{comment.content}</p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
