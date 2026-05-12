# BaseTable Component - Header Customization Guide

## Overview

`BaseTable` kini memiliki fitur header customization yang optional, memberikan flexibility penuh untuk mengatur search, buttons, dan actions.

## Props yang Ditambahkan

```typescript
// Header customization (optional)
searchPlaceholder?: string;          // default: "Search..."
showDeleteButton?: boolean;          // default: true
onDeleteAll?: () => void;            // callback for delete all action
actionButtons?: ReactNode;           // custom React elements untuk buttons
```

## Contoh Penggunaan

### 1. **Default Behavior** (tidak perlu customize apapun)

```typescript
<BaseTable
  data={users}
  columns={columns}
  // ... other required props
  // Header akan tampil dengan default search placeholder dan delete button
/>
```

### 2. **Custom Search Placeholder**

```typescript
<BaseTable
  data={users}
  columns={columns}
  // ... other props
  searchPlaceholder="Find users by name or email..."
/>
```

### 3. **Hide Delete Button**

```typescript
<BaseTable
  data={users}
  columns={columns}
  // ... other props
  showDeleteButton={false}
  // Delete button tidak akan ditampilkan
/>
```

### 4. **Add Delete All Handler**

```typescript
const handleDeleteAllUsers = () => {
  console.log("Delete all users");
  // Implementation
};

<BaseTable
  data={users}
  columns={columns}
  // ... other props
  onDeleteAll={handleDeleteAllUsers}
  showDeleteButton={true}
/>
```

### 5. **Completely Custom Action Buttons**

```typescript
const customButtons = (
  <>
    <Button variant="outline" size="sm" onClick={handleExport}>
      Export
    </Button>
    <Button variant="primary" size="sm" onClick={handleImport}>
      Import
    </Button>
  </>
);

<BaseTable
  data={users}
  columns={columns}
  // ... other props
  actionButtons={customButtons}
  // Ini akan menggantikan default Create + Delete button
/>
```

## Fitur Flexibility

### Component Structure

```
TableHeader (reusable component)
├── Show X entries dropdown
├── Search input
└── Action buttons (Create, Delete All, atau custom)
```

### Keuntungan

✅ Reusable `TableHeader` component untuk berbagai tabel  
✅ Dapat dikustomisasi per tabel tanpa mengubah BaseTable  
✅ Support custom buttons untuk use case spesifik  
✅ Maintain clean separation of concerns  
✅ Backward compatible - default behavior tetap sama

## File Structure

```
components/
├── tables/
│   ├── BaseTable.tsx          ← Main table component (generic, flexible)
│   └── TableHeader.tsx        ← New! Extracted header component
```

## Type Safety

- `ColumnDefinition<T>` tetap fully generic dan type-safe
- Props validation untuk setiap customization option
- TypeScript akan catch errors pada compile time

## Saat Menggunakan Custom Buttons

Jika `actionButtons` prop disediakan:

- Default "Create" button tetap ada
- "Delete All" button digantikan oleh custom buttons
- Pastikan onClick handlers sudah tepat

```typescript
// ❌ Jangan lupa bind event handlers
const customButtons = <Button onClick={() => {}}> Text </Button>;

// ✅ Benar - event handler sudah terdefinisi
const handleCustomAction = () => { /* ... */ };
const customButtons = <Button onClick={handleCustomAction}> Text </Button>;
```

## Best Practice

1. **Untuk list dengan aksi batch** → Gunakan `onDeleteAll`
2. **Untuk list tanpa aksi bulk** → Set `showDeleteButton={false}`
3. **Untuk custom workflow** → Gunakan `actionButtons` prop
4. **Default case** → Kosongkan, gunakan default behavior
