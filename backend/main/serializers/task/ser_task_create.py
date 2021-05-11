from rest_framework import serializers

from .ser_task import TaskSerializer
from ..subtask.ser_subtask import SubtaskSerializer
from ...models import Column, Task
from ...validation.val_auth import \
    authenticate, authorize, authorization_error


class CreateTaskSerializer(TaskSerializer):
    column = serializers.PrimaryKeyRelatedField(
        queryset=Column.objects.select_related('board')
                               .prefetch_related('task_set')
                               .all(),
        error_messages={'does_not_exist': 'Column does not exist.',
                        'null': 'Column cannot be null.'}
    )
    subtasks = serializers.ListField(allow_null=True)
    auth_user = serializers.CharField(allow_blank=True)
    auth_token = serializers.CharField(allow_blank=True)

    class Meta:
        model = TaskSerializer.Meta.model
        fields = \
            'column', 'title', 'description', 'auth_user', 'auth_token', \
            'subtasks'
        extra_kwargs = {
            'title': {
                'error_messages': {
                    'blank': 'Title cannot be blank.',
                    'max_length': 'Title cannot be longer than 50 characters.'
                }
            },
        }

    def validate(self, attrs):
        user = authenticate(attrs.get('auth_user'), attrs.get('auth_token'))
        column = attrs.get('column')
        authorize(user, column.board.team_id)
        return attrs

    def create(self, validated_data):
        column = validated_data.get('column')

        task_serializer = TaskSerializer(
            data={'title': validated_data.get('title'),
                  'description': validated_data.get('description'),
                  'order': 0,
                  'column': column.id}
        )
        task_serializer.is_valid(raise_exception=True)
        task = task_serializer.save()

        subtasks = validated_data.get('subtasks')
        subtasks_data = [
            {'title': subtask, 'order': i, 'task': task.id}
            for i, subtask in enumerate(subtasks)
        ] if subtasks else []

        subtask_serializer = SubtaskSerializer(data=subtasks_data, many=True)
        if not subtask_serializer.is_valid():
            task.delete()
            raise serializers.ValidationError({
                'subtask': subtask_serializer.errors
            })
        subtask_serializer.save()

        for task in column.task_set.all():
            task.order += 1
        Task.objects.bulk_update(column.task_set.all(), ['order'])

        self.instance = task
        return task

    def to_representation(self, instance):
        return {'msg': 'Task creation successful.', 'task_id': instance.id}
