<template>
    <div class="details-layout">
        <div ref="detailsWrapper">
            <slot name="details-header"></slot>
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="$groupedProperties[groupIndex].length">
                    <cmdb-collapse
                        :label="group['bk_group_name']"
                        :collapse.sync="collapseStatus[group['bk_group_id']]">
                        <ul class="property-list clearfix">
                            <li class="property-item clearfix fl"
                                v-for="(property, propertyIndex) in $groupedProperties[groupIndex]"
                                :key="propertyIndex"
                                :title="getTitle(inst, property)">
                                <span class="property-name fl">{{property['bk_property_name']}}</span>
                                <span class="property-value clearfix fl" v-if="property.unit">
                                    <span class="property-value-text fl">{{getValue(property)}}</span>
                                    <span class="property-value-unit fl">{{property.unit}}</span>
                                </span>
                                <span class="property-value fl" v-else>{{getValue(property)}}</span>
                            </li>
                        </ul>
                    </cmdb-collapse>
                </div>
            </template>
        </div>
        <div class="details-options"
            v-if="showOptions"
            :class="{ sticky: scrollbar }">
            <slot name="details-options">
                <span class="inline-block-middle"
                    v-if="showEdit"
                    v-cursor="{
                        active: !$isAuthorized(editAuth),
                        auth: [editAuth]
                    }">
                    <bk-button class="button-edit" type="primary"
                        :disabled="!$isAuthorized(editAuth)"
                        @click="handleEdit">
                        {{editText}}
                    </bk-button>
                </span>
                <span class="inline-block-middle"
                    v-if="showDelete"
                    v-cursor="{
                        active: !$isAuthorized(deleteAuth),
                        auth: [deleteAuth]
                    }">
                    <bk-button class="button-delete" type="danger"
                        :disabled="!$isAuthorized(deleteAuth)"
                        @click="handleDelete">
                        {{deleteText}}
                    </bk-button>
                </span>
            </slot>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
        name: 'cmdb-details',
        mixins: [formMixins],
        props: {
            inst: {
                type: Object,
                required: true
            },
            showOptions: {
                type: Boolean,
                default: true
            },
            editButtonText: {
                type: String,
                default: ''
            },
            deleteButtonText: {
                type: String,
                default: ''
            },
            showEdit: {
                type: Boolean,
                default: true
            },
            showDelete: {
                type: Boolean,
                default: true
            },
            editAuth: {
                type: [String, Array],
                default: ''
            },
            deleteAuth: {
                type: [String, Array],
                default: ''
            }
        },
        data () {
            return {
                collapseStatus: {
                    none: true
                },
                scrollbar: false
            }
        },
        computed: {
            editText () {
                return this.editButtonText || this.$t("Common['编辑']")
            },
            deleteText () {
                return this.deleteButtonText || this.$t("Common['删除']")
            }
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.detailsWrapper, this.checkScrollbar)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$el.detailsWrapper, this.checkScrollbar)
        },
        methods: {
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            handleToggleGroup (group) {
                const groupId = group['bk_group_id']
                const collapse = !!this.collapseStatus[groupId]
                this.$set(this.collapseStatus, groupId, !collapse)
            },
            getTitle (inst, property) {
                return `${property['bk_property_name']}: ${inst[property['bk_property_id']] || '--'} ${property.unit}`
            },
            getValue (property) {
                const value = this.inst[property['bk_property_id']]
                return String(value).length ? value : '--'
            },
            handleEdit () {
                this.$emit('on-edit', this.inst)
            },
            handleDelete () {
                this.$emit('on-delete', this.inst)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-layout{
        height: 100%;
        padding: 0 0 0 32px;
        @include scrollbar-y;
    }
    .property-group{
        padding: 7px 0 10px 0;
        &:first-child{
            padding: 28px 0 10px 0;
        }
    }
    .group-name{
        font-size: 14px;
        line-height: 14px;
        color: #333948;
        overflow: visible;
        .group-toggle {
            cursor: pointer;
            &.collapse .bk-icon {
                transform: rotate(-90deg);
            }
            .bk-icon {
                vertical-align: baseline;
                font-size: 12px;
                font-weight: bold;
                transition: transform .2s ease-in-out;
            }
        }
    }
    .property-list{
        padding: 4px 0;
        .property-item{
            width: 50%;
            max-width: 400px;
            margin: 12px 0 0;
            font-size: 12px;
            line-height: 16px;
            .property-name{
                position: relative;
                width: 35%;
                padding: 0 16px 0 0;
                text-align: right;
                color: $cmdbTextColor;
                @include ellipsis;
                &:after{
                    content: ":";
                    position: absolute;
                    right: 10px;
                }
            }
            .property-value{
                width: 65%;
                padding: 0 15px 0 0;
                @include ellipsis;
                &-text{
                    display: block;
                    max-width: calc(100% - 60px);
                    @include ellipsis;
                }
                &-unit{
                    display: block;
                    width: 60px;
                    padding: 0 0 0 5px;
                    @include ellipsis;
                }
            }
        }
    }
    .details-options{
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 28px 18px 0;
        &.sticky {
            width: calc(100% + 32px);
            margin: 0 0 0 -40px;
            padding: 10px 50px;
            border-top: 1px solid $cmdbBorderColor;
            background-color: #fff;
        }
        .button-edit{
            min-width: 76px;
            margin-right: 4px;
        }
        .button-delete{
            min-width: 76px;
            background-color: #fff;
            color: #ff5656;
        }
    }
</style>
