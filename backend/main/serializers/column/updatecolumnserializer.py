from rest_framework import serializers
import status

from .columnserializer import ColumnSerializer
from ..task.taskserializer import TaskSerializer
from ...models import Column, Task
from ...validation.auth import \
    authenticate_custom, authorize_custom, authorization_error
from ...validation.custom import CustomAPIException


class UpdateColumnSerializer(ColumnSerializer):
    column = serializers.PrimaryKeyRelatedField(
        queryset=Column.objects.select_related('board').all(),
        error_messages={
            'blank': 'Column ID cannot be blank.',
            'null': 'Column ID cannot be null.',
            'invalid': 'Column ID must be a number.'
        },
    )
    tasks = serializers.ListField(allow_empty=True)
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = ColumnSerializer.Meta.model
        fields = 'column', 'tasks', 'auth_user', 'auth_token'

    def validate(self, attrs):
        username = attrs.pop('auth_user')
        token = attrs.pop('auth_token')

        user, authentication_error = authenticate_custom(username, token)
        if authentication_error:
            raise authentication_error
        attrs['user'] = user

        local_authorization_error = authorize_custom(user.username)
        if local_authorization_error:
            # save it for later, as some non-admin users can also use this
            attrs['local_authorization_error'] = local_authorization_error

        column = attrs.pop('column')
        if column.board.team_id != user.team_id:
            raise authorization_error

        self.instance = column
        return attrs

    def update(self, instance, validated_data):
        board_tasks = Task.objects.filter(column__board_id=instance.board_id)
        for task in validated_data.get('tasks'):
            try:
                task_id = task.pop('id')
            except KeyError:
                raise CustomAPIException('task.id',
                                         'Task ID cannot be empty.',
                                         status.HTTP_400_BAD_REQUEST)

            existing_task = board_tasks.get(id=task_id)

            user = validated_data.get('user')
            local_authorization_error = validated_data.get(
                'local_authorization_error'
            )
            if local_authorization_error \
                    and task.get('user') != user.username \
                    and instance.id != existing_task.column_id:
                raise local_authorization_error

            serializer = TaskSerializer(existing_task,
                                        data={**task, 'column': instance.id},
                                        partial=True)
            if serializer.is_valid():
                serializer.save()
            else:
                raise serializer.errors

        return instance

    def to_representation(self, instance):
        return {
            'msg': 'Column and all its tasks updated successfully.',
            'id': instance.id,
        }
