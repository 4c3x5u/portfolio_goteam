from rest_framework import serializers

from .ser_task import TaskSerializer
from ...models import Task
from ...validation.val_auth import authenticate_custom, authorize_custom, \
    authorization_error


class DeleteTaskSerializer(TaskSerializer):
    task = serializers.PrimaryKeyRelatedField(
        queryset=Task.objects.select_related('column', 'column__board').all(),
        error_messages={'null': 'Task ID cannot be null.',
                        'invalid': 'Task ID must be a number.',
                        'incorrect_type': 'Task ID must be a number.',
                        'does_not_exist': 'Task does not exist.'}
    )
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = TaskSerializer.Meta.model
        fields = 'task', 'auth_user', 'auth_token'

    def validate(self, attrs):
        auth_user = attrs.get('auth_user')
        auth_token = attrs.get('auth_token')

        user, authentication_error = authenticate_custom(auth_user, auth_token)
        if authentication_error:
            raise authentication_error

        local_authorization_error = authorize_custom(auth_user)
        if local_authorization_error:
            raise local_authorization_error

        task = attrs.get('task')

        if task.column.board.team_id != user.team_id:
            raise authorization_error

        self.instance = task.id
        return task

    def delete(self):
        return self.validated_data.delete()

    def to_representation(self, instance):
        return {
            'msg': 'Task deleted successfully.',
            'id': instance,
        }