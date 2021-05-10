from rest_framework import serializers
import status

from .taskserializer import TaskSerializer
from ..subtask.subtaskserializer import SubtaskSerializer
from ...models import Task, User
from ...validation.auth import authenticate_custom, authorize_custom, \
    authorization_error
from ...validation.custom import CustomAPIException
from ...validation.column import validate_column_id_custom


class UpdateTaskSerializer(TaskSerializer):
    task = serializers.PrimaryKeyRelatedField(
        queryset=Task.objects.select_related('column', 'column__board')
                             .prefetch_related('subtask_set')
                             .all(),
        error_messages={'blank': 'Task ID cannot be empty.',
                        'invalid': 'Task ID must be a number.'}
    )
    user = serializers.CharField(required=False)
    title = serializers.CharField(
        required=False,
        error_messages={'blank': 'Task title cannot be blank.'}
    )
    order = serializers.IntegerField(
        required=False,
        error_messages={'invalid': 'Task order must be a number.'}
    )
    column = serializers.IntegerField(
        required=False,
        error_messages={'invalid': 'Column ID must be a number.'}
    )
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    def validate(self, attrs):
        auth_user = attrs.pop('auth_user')
        auth_token = attrs.pop('auth_token')

        authenticated_user, authentication_error = \
            authenticate_custom(auth_user, auth_token)
        if authentication_error:
            raise authentication_error

        local_authorization_error = authorize_custom(auth_user)
        if local_authorization_error:
            raise local_authorization_error

        task = attrs.pop('task')
        if task.column.board.team_id != authenticated_user.team_id:
            raise authorization_error

        if 'title' in attrs.keys() and not attrs.get('title'):
            raise CustomAPIException('title',
                                     'Task title cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

        order = attrs.get('order')
        if 'order' in attrs.keys() and (order == '' or order is None):
            raise CustomAPIException('order',
                                     'Task order cannot be empty.',
                                     status.HTTP_400_BAD_REQUEST)

        if 'column' in attrs.keys():
            column_id = attrs.get('column')
            validate_column_id_custom(column_id)

        self.instance = task
        return attrs

    def update(self, instance, validated_data):
        subtasks = validated_data.pop(
            'subtasks'
        ) if 'subtasks' in validated_data.keys() else None

        # update task
        task_serializer = TaskSerializer(instance,
                                         data=validated_data,
                                         partial=True)
        if not task_serializer.is_valid():
            raise serializers.ValidationError({'task': task_serializer.errors})
        task = task_serializer.save()

        # update subtasks
        if subtasks:
            task.subtask_set.all().delete()
            subtasks_data = [
                {
                    'title': subtask['title'],
                    'order': subtask['order'],
                    'task': task.id,
                    'done': subtask['done']
                } for subtask in subtasks
            ]
            subtask_serializer = SubtaskSerializer(data=subtasks_data,
                                                   many=True)
            if not subtask_serializer.is_valid():
                raise serializers.ValidationError(
                    {'subtasks': subtask_serializer.errors}
                )
            subtask_serializer.save()

        return task

    def to_representation(self, instance):
        return {
            'msg': 'Task update successful.',
            'id': instance.id
        }
