from rest_framework import serializers

from .ser_task import TaskSerializer
from ..subtask.ser_subtask import SubtaskSerializer
from ...models import Task
from ...helpers.auth_helper import AuthHelper


class UpdateTaskSerializer(TaskSerializer):
    task = serializers.PrimaryKeyRelatedField(
        queryset=Task.objects.select_related('column', 'column__board')
                             .prefetch_related('subtask_set')
                             .all(),
        error_messages={'blank': 'Task ID cannot be empty.',
                        'invalid': 'Task ID must be a number.'}
    )
    subtasks = serializers.ListField(allow_null=True)
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

    class Meta(TaskSerializer.Meta):
        pass

    def validate(self, attrs):
        user = AuthHelper.authenticate(attrs.pop('auth_user'),
                                       attrs.pop('auth_token'))

        task = attrs.pop('task')
        AuthHelper.authorize(user, task.column.board.team_id)

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
        task.subtask_set.all().delete()
        if subtasks:
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
