---
aside: false
outline: false
---

<script setup lang="ts">
import { useData } from 'vitepress'
import { OAOperation } from 'vitepress-openapi/client'

const { page } = useData()
const operationId = page.value.params?.operationId
</script>

<OAOperation :operationId="operationId" />
