'use client';

import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { useCategories, useDeleteCategory } from '@/lib/hooks/use-api';
import { useAuth } from '@/lib/hooks/use-auth';
import { Edit, Trash2 } from 'lucide-react';

export default function CategoriesPage() {
  const { currentStore } = useAuth();
  const storeId = currentStore?.id ?? '';
  const { data: categories, isLoading } = useCategories(storeId);
  const deleteCategoryMutation = useDeleteCategory();

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this category?')) return;
    try {
      await deleteCategoryMutation.mutateAsync({ storeId, categoryId: id });
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Categories</h1>
          <p className="text-gray-600">Manage your product categories.</p>
        </div>
        <Button asChild>
          <Link href="/dashboard/categories/new">Create Category</Link>
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Categories</CardTitle>
          <CardDescription>
            Categories help organize your products.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Slug</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {categories?.map((category) => (
                  <TableRow key={category.id}>
                    <TableCell className="text-sm font-medium text-gray-900">{category.name}</TableCell>
                    <TableCell className="text-sm text-gray-600">{category.slug}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Button variant="outline" size="sm" asChild>
                          <Link href={`/dashboard/categories/${category.id}`}>
                            <Edit className="mr-2 h-4 w-4" />
                            Edit
                          </Link>
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="text-red-600"
                          onClick={() => handleDelete(category.id)}
                          disabled={deleteCategoryMutation.isPending}
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          Delete
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
                {(!categories || categories.length === 0) && (
                  <TableRow>
                    <TableCell colSpan={3} className="h-24 text-center">
                      {isLoading ? 'Loading categories...' : 'No categories found.'}
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
